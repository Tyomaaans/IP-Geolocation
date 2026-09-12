package ip

import (
    "context"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/google/uuid"
    amqp "github.com/rabbitmq/amqp091-go"
    "github.com/redis/go-redis/v9"

	"ip-geo/internal/domain"
    "ip-geo/pkg"
)

const (
    dailyIpLimit    = 900
    redisCounterKey = "ip:daily:count"
    queueName       = "ip.queue"
    retryQueueName  = "ip.retry.queue"
    retryExchange   = "ip.retry.exchange"
    mainExchange    = "ip.main.exchange"
)

type IpService interface {
    TrackVisit(ctx context.Context, ip, ua string) error
    DeliverQueued(ctx context.Context, queueIp domain.IpQueueEntity) error

    GetIpVisitByID(ctx context.Context, rawID string) (*IpResponse, error)
    GetIpVisitByIP(ctx context.Context, ip string) (*IpResponse, error)
    GetIpVisits(ctx context.Context, page, limit int) (*PaginatedIpVisitResponse, error)
    GetTodayIpVisits(ctx context.Context, page, limit int) (*PaginatedIpVisitResponse, error)

    GetIpHistoryByID(ctx context.Context, rawID string) (*IpResponse, error)
    GetIpHistoriesByIP(ctx context.Context, ip string, page, limit int) (*PaginatedIpHistoryResponse, error)
    GetIpHistories(ctx context.Context, page, limit int) (*PaginatedIpHistoryResponse, error)
    GetTodayIpHistories(ctx context.Context, page, limit int) (*PaginatedIpHistoryResponse, error)
}

type ipService struct {
    ipRepo      domain.IpRepository
    ipGeoClient domain.IpGeoClient
    redis       *redis.Client
    amqpCh      *amqp.Channel
}

func NewIpService(
    ipRepo      domain.IpRepository,
    ipGeoClient domain.IpGeoClient,
    redisClient *redis.Client,
    amqpCh      *amqp.Channel,
) IpService {
    return &ipService{
        ipRepo:      ipRepo,
        ipGeoClient: ipGeoClient,
        redis:       redisClient,
        amqpCh:      amqpCh,
    }
}

func (s *ipService) TrackVisit(ctx context.Context, ip, ua string) error {
    existing, err := s.ipRepo.GetIpVisitByIP(ctx, ip)
    if err != nil && !errors.Is(err, pkg.ErrNotFound) {
        return err
    }

    id := uuid.NewString()

    if existing != nil {
        history := ToIpHistoryEntityFromIpVisitEntity(id, ua, existing)
        if err := s.ipRepo.CreateIpHistory(ctx, history); err != nil {
            return err
        }
        return nil
    }

    // Buat queue entity untuk diproses via RabbitMQ & Redis limit check
    queueIp := domain.IpQueueEntity{
        ID:        id,
        IpAddress: ip,
        UserAgent: ua,
        Attempt:   0,
    }

    return s.enqueue(ctx, queueIp)
}

func (s *ipService) DeliverQueued(ctx context.Context, queueIp domain.IpQueueEntity) error {
    allowed, err := s.checkAndIncrement(ctx)
    if err != nil {
        return err
    }

    if !allowed {
        log.Printf("ip-geo: daily limit reached! retrying for IP %s", queueIp.IpAddress)
        return s.enqueueRetry(ctx, queueIp)
    }

    return s.processVisit(ctx, queueIp)
}

// Internal Queue & Redis Helpers

func (s *ipService) enqueue(ctx context.Context, queueIp domain.IpQueueEntity) error {
    body, err := json.Marshal(queueIp)
    if err != nil {
        return err
    }

    // Optional: Simpan status pending/tracking awal ke DB jika diperlukan repository-nya
    // Seperti pada email service yang mencatat entity ke database sebelum di-publish ke broker.

    err = s.amqpCh.PublishWithContext(ctx, mainExchange, queueName, false, false, amqp.Publishing{
        ContentType:  "application/json",
        DeliveryMode: amqp.Persistent,
        Priority:     0,
        Body:         body,
    })
    if err != nil {
        return pkg.HandleRabbitError(err)
    }

    return nil
}

func (s *ipService) enqueueRetry(ctx context.Context, queueIp domain.IpQueueEntity) error {
    body, err := json.Marshal(queueIp)
    if err != nil {
        return err
    }

    queueIp.Attempt++
    delay := s.delayUntilNextReset()

    err = s.amqpCh.PublishWithContext(ctx, retryExchange, retryQueueName, false, false, amqp.Publishing{
        ContentType:  "application/json",
        DeliveryMode: amqp.Persistent,
        Priority:     10,
        Expiration:   fmt.Sprintf("%d", delay.Milliseconds()),
        Body:         body,
    })
    if err != nil {
        return pkg.HandleRabbitError(err)
    }

    return nil
}

func (s *ipService) delayUntilNextReset() time.Duration {
    now := time.Now().UTC()
    nextReset := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
    return time.Until(nextReset)
}

func (s *ipService) processVisit(ctx context.Context, queueIp domain.IpQueueEntity) error {
    geoData, err := s.ipGeoClient.FetchByIP(ctx, queueIp.IpAddress)
    if err != nil {
        return err
    }

    visit := ToIpVisitEntityFromIpGeo(queueIp.ID, queueIp.UserAgent, geoData)
    if err := s.ipRepo.CreateIpVisit(ctx, visit); err != nil {
        return err
    }

    return nil
}

func (s *ipService) checkAndIncrement(ctx context.Context) (bool, error) {
    pipe := s.redis.TxPipeline()

    incr := pipe.Incr(ctx, redisCounterKey)

    now := time.Now().UTC()
    midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
    ttl := time.Until(midnight)
    pipe.ExpireNX(ctx, redisCounterKey, ttl)

    if _, err := pipe.Exec(ctx); err != nil {
        return false, pkg.HandleRedisError(err)
    }

    count := incr.Val()
    if count > int64(dailyIpLimit) {
        if err := s.redis.Decr(ctx, redisCounterKey).Err(); err != nil {
            return false, pkg.HandleRedisError(err)
        }
        return false, nil
    }

    return true, nil
}

// IP Visit GET Methods

func (s *ipService) GetIpVisitByID(ctx context.Context, rawID string) (*IpResponse, error) {
    id, err := parseOrDecodeUUID(rawID)
    if err != nil {
        return nil, err
    }

    visit, err := s.ipRepo.GetIpVisitByID(ctx, id)
    if err != nil {
        return nil, err
    }

    res, err := ToIpVisitResponse(visit)
    if err != nil {
        return nil, err
    }
    
    return res, nil
}

func (s *ipService) GetIpVisitByIP(ctx context.Context, ip string) (*IpResponse, error) {
    visit, err := s.ipRepo.GetIpVisitByIP(ctx, ip)
    if err != nil {
        return nil, err
    }

    res, err := ToIpVisitResponse(visit)
    if err != nil {
        return nil, err
    }
    
    return res, nil
}

func (s *ipService) GetIpVisits(ctx context.Context, page, limit int) (*PaginatedIpVisitResponse, error) {
    visits, pages, err := s.ipRepo.GetIpVisits(ctx, page, limit)
    if err != nil {
        return nil, err
    }

    res, err := ToPaginatedIpVisitResponse(visits, page, limit, pages)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (s *ipService) GetTodayIpVisits(ctx context.Context, page, limit int) (*PaginatedIpVisitResponse, error) {
    visits, pages, err := s.ipRepo.GetTodayIpVisits(ctx, page, limit)
    if err != nil {
        return nil, err
    }

    res, err := ToPaginatedIpVisitResponse(visits, page, limit, pages)
    if err != nil {
        return nil, err
    }

    return res, nil
}

// IP History GET Methods

func (s *ipService) GetIpHistoryByID(ctx context.Context, rawID string) (*IpResponse, error) {
    id, err := parseOrDecodeUUID(rawID)
    if err != nil {
        return nil, err
    }
    
    history, err := s.ipRepo.GetIpHistoryByID(ctx, id)
    if err != nil {
        return nil, err
    }

    res, err := ToIpHistoryResponse(history)
    if err != nil {
        return nil, err
    }
    
    return res, nil
}

func (s *ipService) GetIpHistoriesByIP(ctx context.Context, ip string, page, limit int) (*PaginatedIpHistoryResponse, error) {
    histories, pages, err := s.ipRepo.GetIpHistoriesByIP(ctx, ip, page, limit)
    if err != nil {
        return nil, err
    }

    res, err := ToPaginatedIpHistoryResponse(histories, page, limit, pages)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (s *ipService) GetIpHistories(ctx context.Context, page, limit int) (*PaginatedIpHistoryResponse, error) {
    histories, pages, err := s.ipRepo.GetIpHistories(ctx, page, limit)
    if err != nil {
        return nil, err
    }

    res, err := ToPaginatedIpHistoryResponse(histories, page, limit, pages)
    if err != nil {
        return nil, err
    }

    return res, nil
}

func (s *ipService) GetTodayIpHistories(ctx context.Context, page, limit int) (*PaginatedIpHistoryResponse, error) {
    histories, pages, err := s.ipRepo.GetTodayIpHistories(ctx, page, limit)
    if err != nil {
        return nil, err
    }

    res, err := ToPaginatedIpHistoryResponse(histories, page, limit, pages)
    if err != nil {
        return nil, err
    }

    return res, nil
}

// Internal Helper

func parseOrDecodeUUID(input string) (string, error) {
    if input == "" {
        return "", pkg.ErrInvalidInput
    }

    if _, err := uuid.Parse(input); err == nil {
        return input, nil
    }

    decoded, err := base64.RawURLEncoding.DecodeString(input)
    if err != nil {
        return "", pkg.ErrInvalidInput
    }

    parsed, err := uuid.FromBytes(decoded)
    if err != nil {
        return "", pkg.ErrInvalidInput
    }

    return parsed.String(), nil
}