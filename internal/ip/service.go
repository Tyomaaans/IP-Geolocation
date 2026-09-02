package ip

import (
	"context"
	"encoding/base64"
	"errors"

	"ip-geo/internal/domain"
	"ip-geo/pkg"

	"github.com/google/uuid"
)

type IpService interface {
	TrackVisit(ctx context.Context, ip, ua string) error

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
}

func NewIpService(
	ipRepo      domain.IpRepository,
	ipGeoClient domain.IpGeoClient,
) IpService {
	return &ipService{
		ipRepo:      ipRepo,
		ipGeoClient: ipGeoClient,
	}
}

func (s *ipService) TrackVisit(ctx context.Context, ip, ua string) error {
	existing, err := s.ipRepo.GetIpVisitByIP(ctx, ip)
	if err != nil && !errors.Is(err, pkg.ErrNotFound){
		return err
	}

	if existing != nil {
		id := uuid.NewString()
		history := ToIpHistoryEntityFromIpVisitEntity(id, ua, existing)
		if err := s.ipRepo.CreateIpHistory(ctx, history); err != nil {
			return err
		}
		return nil
	}

	geoData, err := s.ipGeoClient.FetchByIP(ctx, ip)
	if err != nil {
		return err
	}

	id := uuid.NewString()

	visit := ToIpVisitEntityFromIpGeo(id, ua, geoData)
	if err := s.ipRepo.CreateIpVisit(ctx, visit); err != nil {
		return err
	}

	return nil
}

// IP Visit

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

// IP History

func (s *ipService) GetIpHistoryByID(ctx context.Context, rawID string) (*IpResponse, error) {
	id, err := parseOrDecodeUUID(rawID)
    if err != nil {
        return nil, err
    }
	
	History, err := s.ipRepo.GetIpHistoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res, err := ToIpHistoryResponse(History)
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