package sqlite

import (
	"time"
)

type IpHistoryStorage struct {
    ID        string    `gorm:"primaryKey"`
    IpAddress string    `gorm:"type:varchar(255);not null;index"`
    UserAgent string    `gorm:"type:varchar(255);not null"`
    Continent string    `gorm:"type:varchar(255);not null"`
    Country   string    `gorm:"type:varchar(255);not null"`
    Province  string    `gorm:"type:varchar(255);not null"`
    District  string    `gorm:"type:varchar(255);not null"`
    City      string    `gorm:"type:varchar(255);not null"`
    ZipCode   string    `gorm:"type:varchar(255);not null"`
    Latitude  string    `gorm:"type:varchar(255);not null"`
    Longitude string    `gorm:"type:varchar(255);not null"`
    CreatedAt time.Time `gorm:"not null"`
}