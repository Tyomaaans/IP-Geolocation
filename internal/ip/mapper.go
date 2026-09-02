package ip

import (
	"encoding/base64"

	"github.com/google/uuid"

	"ip-geo/internal/domain"
	"ip-geo/internal/infrastructure/sqlite"
)

// IP Visit to Storage <-> Entity

func ToIpVisitEntity(s *sqlite.IpVisitStorage) *domain.IpVisitEntity {
	return &domain.IpVisitEntity{
		ID:        s.ID,
		IpAddress: s.IpAddress,
		UserAgent: s.UserAgent,
		Continent: s.Continent,
		Country:   s.Country,
		Province:  s.Province,
		District:  s.District,
		City:      s.City,
		ZipCode:   s.ZipCode,
		Latitude:  s.Latitude,
		Longitude: s.Longitude,
		VisitTime: s.CreatedAt,
	}
}

func ToIpVisitStorage(e *domain.IpVisitEntity) *sqlite.IpVisitStorage {
	return &sqlite.IpVisitStorage{
		ID:         e.ID,
		IpAddress:  e.IpAddress,
		UserAgent:  e.UserAgent,
		Continent:  e.Continent,
		Country:    e.Country,
		Province:   e.Province,
		District:   e.District,
		City:       e.City,
		ZipCode:    e.ZipCode,
		Latitude:   e.Latitude,
		Longitude: e.Longitude,
	}
}

func ToIpVisitListEntity(list []sqlite.IpVisitStorage) []domain.IpVisitEntity {
	result := make([]domain.IpVisitEntity, len(list))

	for i := range list {
		result[i] = *ToIpVisitEntity(&list[i])
	}

	return result
}

func ToIpVisitEntityFromIpGeo(id, ua string, e *domain.IpGeoEntity) *domain.IpVisitEntity {
	return &domain.IpVisitEntity{
		ID:        id,
		IpAddress: e.IP,
		UserAgent: ua,
		Continent: e.Location.ContinentName,
		Country:   e.Location.CountryName,
		Province:  e.Location.StateProv,
		District:  e.Location.District,
		City:      e.Location.City,
		ZipCode:   e.Location.Zipcode,
		Latitude:  e.Location.Latitude,
		Longitude: e.Location.Longitude,
	}
}

// IP History to Storage <-> Entity

func ToIpHistoryEntity(s *sqlite.IpHistoryStorage) *domain.IpHistoryEntity {
	return &domain.IpHistoryEntity{
		ID:         s.ID,
		IpAddress:  s.IpAddress,
		UserAgent:  s.UserAgent,
		Continent:  s.Continent,
		Country:    s.Country,
		Province:   s.Province,
		District:   s.District,
		City:       s.City,
		ZipCode:    s.ZipCode,
		Latitude:   s.Latitude,
		Longitude: s.Longitude,
		VisitTime:  s.CreatedAt,
	}
}

func ToIpHistoryStorage(e *domain.IpHistoryEntity) *sqlite.IpHistoryStorage {
	return &sqlite.IpHistoryStorage{
		ID:         e.ID,
		IpAddress:  e.IpAddress,
		UserAgent:  e.UserAgent,
		Continent:  e.Continent,
		Country:    e.Country,
		Province:   e.Province,
		District:   e.District,
		City:       e.City,
		ZipCode:    e.ZipCode,
		Latitude:   e.Latitude,
		Longitude: e.Longitude,
	}
}

func ToIpHistoryListEntity(list []sqlite.IpHistoryStorage) []domain.IpHistoryEntity {
	result := make([]domain.IpHistoryEntity, len(list))

	for i := range list {
		result[i] = *ToIpHistoryEntity(&list[i])
	}

	return result
}

// IP Visit <-> History

func ToIpVisitEntityFromIpHistroryEntity(e *domain.IpHistoryEntity) *domain.IpVisitEntity {
	return &domain.IpVisitEntity{
		ID:         e.ID,
		IpAddress:  e.IpAddress,
		UserAgent:  e.UserAgent,
		Continent:  e.Continent,
		Country:    e.Country,
		Province:   e.Province,
		District:   e.District,
		City:       e.City,
		ZipCode:    e.ZipCode,
		Latitude:   e.Latitude,
		Longitude: e.Longitude,
		VisitTime:  e.VisitTime,
	}
}

func ToIpHistoryEntityFromIpVisitEntity(id, ua string, e *domain.IpVisitEntity) *domain.IpHistoryEntity {
	return &domain.IpHistoryEntity{
		ID:        id, 
		IpAddress: e.IpAddress,
		UserAgent: ua,
		Continent: e.Continent,
		Country:   e.Country,
		Province:  e.Province,
		District:  e.District,
		City:      e.City,
		ZipCode:   e.ZipCode,
		Latitude:  e.Latitude,
		Longitude: e.Longitude,
		VisitTime: e.VisitTime,
	}
}

// IP Response

func ToIpVisitResponse(e *domain.IpVisitEntity) (*IpResponse, error) {
	id, err := uuidToBase64(e.ID)
	if err != nil {
		return nil, err
	}

	return &IpResponse{
		ID:        id,
		IpAddress: e.IpAddress,
		UserAgent: e.UserAgent,
		Continent: e.Continent,
		Country:   e.Country,
		Province:  e.Province,
		District:  e.District,
		City:      e.City,
		ZipCode:   e.ZipCode,
		Latitude:  e.Latitude,
		Longitude: e.Longitude,
		VisitTime: e.VisitTime,
	}, nil
}

func ToIpVisitListResponse(list []domain.IpVisitEntity) ([]IpResponse, error) {
    result := make([]IpResponse, len(list))

    for i := range list {
        res, err := ToIpVisitResponse(&list[i])
        if err != nil {
            return nil, err
        }
        result[i] = *res
    }

    return result, nil
}

func ToPaginatedIpVisitResponse(list []domain.IpVisitEntity, page, limit int, totalData int64) (*PaginatedIpVisitResponse, error) {
    listResponse, err := ToIpVisitListResponse(list)
    if err != nil {
        return nil, err
    }

    if limit <= 0 {
        limit = 10
    }

    totalPages := int((totalData + int64(limit) - 1) / int64(limit))

    return &PaginatedIpVisitResponse{
        Data: listResponse,
        Meta: PaginationMeta{
            CurrentPage: page,
            TotalPages:  totalPages,
            PageSize:    limit,
            TotalData:   totalData,
        },
    }, nil
}

func ToIpHistoryResponse(e *domain.IpHistoryEntity) (*IpResponse, error) {
	id, err := uuidToBase64(e.ID)
	if err != nil {
		return nil, err
	}

	return &IpResponse{
		ID:        id,
		IpAddress: e.IpAddress,
		UserAgent: e.UserAgent,
		Continent: e.Continent,
		Country:   e.Country,
		Province:  e.Province,
		District:  e.District,
		City:      e.City,
		ZipCode:   e.ZipCode,
		Latitude:  e.Latitude,
		Longitude: e.Longitude,
		VisitTime: e.VisitTime,
	}, nil
}

func ToIpHistoryListResponse(list []domain.IpHistoryEntity) ([]IpResponse, error) {
    result := make([]IpResponse, len(list))

    for i := range list {
        res, err := ToIpHistoryResponse(&list[i])
        if err != nil {
            return nil, err
        }
        result[i] = *res
    }

    return result, nil
}

func ToPaginatedIpHistoryResponse(list []domain.IpHistoryEntity, page, limit int, totalData int64) (*PaginatedIpHistoryResponse, error) {
    listResponse, err := ToIpHistoryListResponse(list)
    if err != nil {
        return nil, err
    }

    if limit <= 0 {
        limit = 10
    }

    totalPages := int((totalData + int64(limit) - 1) / int64(limit))

    return &PaginatedIpHistoryResponse{
        Data: listResponse,
        Meta: PaginationMeta{
            CurrentPage: page,
            TotalPages:  totalPages,
            PageSize:    limit,
            TotalData:   totalData,
        },
    }, nil
}

func uuidToBase64(uuidStr string) (string, error) {
	parsed, err := uuid.Parse(uuidStr)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(parsed[:]), nil
}