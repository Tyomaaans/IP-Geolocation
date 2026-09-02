package domain

import (
	"time"
)

type IpVisitEntity struct {
	ID        string
	IpAddress string
	UserAgent string
	Continent string
	Country   string
	Province  string
	District  string
	City      string
	ZipCode   string
	Latitude  string
	Longitude string
	VisitTime time.Time
}