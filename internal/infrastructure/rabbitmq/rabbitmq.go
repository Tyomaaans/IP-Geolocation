package rabbitmq

import (
    "context"
    "encoding/json"
    "fmt"
    "log/slog"

    amqp "github.com/rabbitmq/amqp091-go"

    "ip-geo/internal/domain"
    "ip-geo/internal/ip"
)

const (
    ipMainExchange   = "ip.main.exchange"
    ipRetryExchange  = "ip.retry.exchange"
    ipQueueName      = "ip.queue"
    ipRetryQueueName = "ip.retry.queue"
)

type RabbitMQ struct {
    conn   *amqp.Connection
    ch     *amqp.Channel
    logger *slog.Logger
}

func NewRabbitMQClient(dsn string, logger *slog.Logger) (*RabbitMQ, error) {
    conn, err := amqp.Dial(dsn)
    if err != nil {
        return nil, fmt.Errorf("rabbitmq: dial ip %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("rabbitmq: ip channel %w", err)
    }

    r := &RabbitMQ{conn: conn, ch: ch, logger: logger}
    if err := r.setup(); err != nil {
        r.Close()
        return nil, err
    }

    return r, nil
}

func (r *RabbitMQ) Channel() *amqp.Channel { 
    return r.ch 
}

func (r *RabbitMQ) Close() {
    if r.ch != nil {
        r.ch.Close()
    }
    if r.conn != nil {
        r.conn.Close()
    }
}

func (r *RabbitMQ) setup() error {
    if err := r.ch.ExchangeDeclare(ipMainExchange, "direct", true, false, false, false, nil); err != nil {
        return fmt.Errorf("rabbitmq: declare ip main exchange %w", err)
    }

    if err := r.ch.ExchangeDeclare(ipRetryExchange, "direct", true, false, false, false, nil); err != nil {
        return fmt.Errorf("rabbitmq: declare ip retry exchange %w", err)
    }

    if _, err := r.ch.QueueDeclare(ipQueueName, true, false, false, false, nil); err != nil {
        return fmt.Errorf("rabbitmq: declare ip main queue %w", err)
    }
    if err := r.ch.QueueBind(ipQueueName, ipQueueName, ipMainExchange, false, nil); err != nil {
        return fmt.Errorf("rabbitmq: bind ip main queue %w", err)
    }

    retryArgs := amqp.Table{
        "x-dead-letter-exchange":    ipMainExchange,
        "x-dead-letter-routing-key": ipQueueName,
    }
    if _, err := r.ch.QueueDeclare(ipRetryQueueName, true, false, false, false, retryArgs); err != nil {
        return fmt.Errorf("rabbitmq: declare ip retry queue %w", err)
    }
    if err := r.ch.QueueBind(ipRetryQueueName, ipRetryQueueName, ipRetryExchange, false, nil); err != nil {
        return fmt.Errorf("rabbitmq: bind ip retry queue %w", err)
    }

    return nil
}

func (r *RabbitMQ) StartWorker(ctx context.Context, svc ip.IpService) error {
    if err := r.ch.Qos(1, 0, false); err != nil {
        return fmt.Errorf("rabbitmq: ip qos %w", err)
    }

    msgs, err := r.ch.Consume(ipQueueName, "", false, false, false, false, nil)
    if err != nil {
        return fmt.Errorf("rabbitmq: ip consume %w", err)
    }

    r.logger.Info("rabbitmq: ip worker started", slog.String("queue", ipQueueName))

    for {
        select {
        case <-ctx.Done():
            r.logger.Info("rabbitmq: ip worker stopping")
            return nil

        case msg, ok := <-msgs:
            if !ok {
                r.logger.Warn("rabbitmq: ip channel closed")
                return fmt.Errorf("rabbitmq: ip channel closed unexpectedly")
            }

            var payload domain.IpQueueEntity
            if err := json.Unmarshal(msg.Body, &payload); err != nil {
                r.logger.Error("rabbitmq: bad ip message, discarding", slog.Any("error", err))
                msg.Nack(false, false)
                continue
            }

            if err := svc.DeliverQueued(ctx, payload); err != nil {
                r.logger.Error("rabbitmq: failed to process ip visit, nacking",
                    slog.String("ip", payload.IpAddress),
                    slog.Any("error", err),
                )
                msg.Nack(false, false)
                continue
            }

            r.logger.Info("rabbitmq: ip visit processed successfully",
                slog.String("ip", payload.IpAddress),
            )

            msg.Ack(false)
        }
    }
}