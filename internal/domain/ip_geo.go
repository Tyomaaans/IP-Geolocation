package domain

import (
	"context"
)

type IpGeoEntity struct {
	IP       string `json:"ip"`
	Location struct {
		ContinentName string `json:"continent_name"`
		CountryName   string `json:"country_name"`
		StateProv     string `json:"state_prov"`
		District      string `json:"district"`
		City          string `json:"city"`
		Zipcode       string `json:"zipcode"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
	} `json:"location"`
}

type IpGeoClient interface {
	FetchByIP(ctx context.Context, ip string) (*IpGeoEntity, error)
}