package ip

import (
	"context"
	"time"

	"gorm.io/gorm"

	"ip-geo/internal/domain"
	"ip-geo/internal/infrastructure/sqlite"
	"ip-geo/pkg"
)

type ipRepository struct {
	db *gorm.DB
}

func NewIpRepository(db *gorm.DB) domain.IpRepository {
	return &ipRepository{
		db: db,
	}
}

// IP Visit

func (r *ipRepository) CreateIpVisit(ctx context.Context, visit *domain.IpVisitEntity) error {
	storage := ToIpVisitStorage(visit)

	if err := r.db.WithContext(ctx).
		Create(storage).
		Error; err != nil {	
			return pkg.HandleDBError(err)
		}
	
	return nil
}

func (r *ipRepository) GetIpVisitByID(ctx context.Context, id string) (*domain.IpVisitEntity, error) {
	var storage sqlite.IpVisitStorage

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&storage).Error

	if err != nil {
		return nil, pkg.HandleDBError(err)
	}

	return ToIpVisitEntity(&storage), nil
}

func (r *ipRepository) GetIpVisitByIP(ctx context.Context, ip string) (*domain.IpVisitEntity, error) {
	var storage sqlite.IpVisitStorage

	err := r.db.WithContext(ctx).
		Where("ip_address = ?", ip).
		First(&storage).Error

	if err != nil {
		return nil, pkg.HandleDBError(err)
	}

	return ToIpVisitEntity(&storage), nil
}

func (r *ipRepository) GetIpVisits(ctx context.Context, page, limit int) ([]domain.IpVisitEntity, int64, error) {
    var storages []sqlite.IpVisitStorage
    var totalData int64

    if err := r.db.WithContext(ctx).Model(&sqlite.IpVisitStorage{}).Count(&totalData).Error; err != nil {
        return nil, 0, pkg.HandleDBError(err)
    }

    if limit <= 0 {
        limit = 10
    }
    if page <= 0 {
        page = 1
    }

    offset := (page - 1) * limit

    if err := r.db.WithContext(ctx).
        Order("created_at ASC").
        Limit(limit).
        Offset(offset).
        Find(&storages).Error; err != nil {
        return nil, 0, err
    }

    return ToIpVisitListEntity(storages), totalData, nil
}

func (r *ipRepository) GetTodayIpVisits(ctx context.Context, page, limit int) ([]domain.IpVisitEntity, int64, error) {
	var storages []sqlite.IpVisitStorage
	var totalData int64

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-1 * time.Nanosecond)

	query := r.db.WithContext(ctx).
		Model(&sqlite.IpVisitStorage{}).
		Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay)

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, pkg.HandleDBError(err)
	}

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if err := query.
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&storages).Error; err != nil {
		return nil, 0, pkg.HandleDBError(err)
	}

	return ToIpVisitListEntity(storages), totalData, nil
}

// IP History

func (r *ipRepository) CreateIpHistory(ctx context.Context, history *domain.IpHistoryEntity) error {
	storage := ToIpHistoryStorage(history)

	if err := r.db.WithContext(ctx).
		Create(storage).
		Error; err != nil {
			return pkg.HandleDBError(err)
		}
	
	return nil
}

func (r *ipRepository) GetIpHistoryByID(ctx context.Context, id string) (*domain.IpHistoryEntity, error) {
	var storage sqlite.IpHistoryStorage

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&storage).Error

	if err != nil {
		return nil, pkg.HandleDBError(err)
	}

	return ToIpHistoryEntity(&storage), nil
}

func (r *ipRepository) GetIpHistoriesByIP(ctx context.Context, ip string, page, limit int) ([]domain.IpHistoryEntity, int64, error) {
    var storages []sqlite.IpHistoryStorage
    var totalData int64

    if err := r.db.WithContext(ctx).Model(&sqlite.IpHistoryStorage{}).Count(&totalData).Error; err != nil {
        return nil, 0, pkg.HandleDBError(err)
    }

    if limit <= 0 {
        limit = 10
    }
    if page <= 0 {
        page = 1
    }

    offset := (page - 1) * limit

    if err := r.db.WithContext(ctx).
		Where("ip_address = ?", ip).
        Order("created_at ASC").
        Limit(limit).
        Offset(offset).
        Find(&storages).Error; err != nil {
        return nil, 0, err
    }

    return ToIpHistoryListEntity(storages), totalData, nil
}

func (r *ipRepository) GetIpHistories(ctx context.Context, page, limit int) ([]domain.IpHistoryEntity, int64, error) {
    var storages []sqlite.IpHistoryStorage
    var totalData int64

    if err := r.db.WithContext(ctx).Model(&sqlite.IpHistoryStorage{}).Count(&totalData).Error; err != nil {
        return nil, 0, pkg.HandleDBError(err)
    }

    if limit <= 0 {
        limit = 10
    }
    if page <= 0 {
        page = 1
    }

    offset := (page - 1) * limit

    if err := r.db.WithContext(ctx).
        Order("created_at ASC").
        Limit(limit).
        Offset(offset).
        Find(&storages).Error; err != nil {
        return nil, 0, pkg.HandleDBError(err)
    }

    return ToIpHistoryListEntity(storages), totalData, nil
}

func (r *ipRepository) GetTodayIpHistories(ctx context.Context, page, limit int) ([]domain.IpHistoryEntity, int64, error) {
	var storages []sqlite.IpHistoryStorage
	var totalData int64

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-1 * time.Nanosecond)

	query := r.db.WithContext(ctx).
		Model(&sqlite.IpHistoryStorage{}).
		Where("created_at BETWEEN ? AND ?", startOfDay, endOfDay)

	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, pkg.HandleDBError(err)
	}

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if err := query.
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&storages).Error; err != nil {
		return nil, 0, pkg.HandleDBError(err)
	}

	return ToIpHistoryListEntity(storages), totalData, nil
}