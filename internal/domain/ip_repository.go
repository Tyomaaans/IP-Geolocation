package domain

import (
	"context"
)

type IpRepository interface {
	// IP Visit
	CreateIpVisit(ctx context.Context, visit *IpVisitEntity) error
	GetIpVisitByID(ctx context.Context, id string) (*IpVisitEntity, error)
	GetIpVisitByIP(ctx context.Context, ip string) (*IpVisitEntity, error)
	GetIpVisits(ctx context.Context, page, limit int) ([]IpVisitEntity, int64, error)
	GetTodayIpVisits(ctx context.Context, page, limit int) ([]IpVisitEntity, int64, error) 
	
	// IP History
	CreateIpHistory(ctx context.Context, history *IpHistoryEntity) error
	GetIpHistoryByID(ctx context.Context, id string) (*IpHistoryEntity, error)
	GetIpHistoriesByIP(ctx context.Context, ip string, page, limit int) ([]IpHistoryEntity, int64, error)
	GetIpHistories(ctx context.Context, page, limit int) ([]IpHistoryEntity, int64, error)
	GetTodayIpHistories(ctx context.Context, page, limit int) ([]IpHistoryEntity, int64, error)
}
