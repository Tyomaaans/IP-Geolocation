package ip

import (
	"time"
)

type IpResponse struct {
	ID        string    `json:"id"`
	IpAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Continent string    `json:"continent"`
	Country   string    `json:"country"`
	Province  string    `json:"province"`
	District  string    `json:"district"`
	City      string    `json:"city"`
	ZipCode   string    `json:"zip_code"`
	Latitude  string    `json:"latitude"`
	Longitude string    `json:"longitude"`
	VisitTime time.Time `json:"visit_time"`
}

type PaginationMeta struct {
    CurrentPage int   `json:"current_page"`
    TotalPages  int   `json:"total_pages"`
    PageSize    int   `json:"page_size"`
    TotalData   int64 `json:"total_data"`
}

type PaginatedIpVisitResponse struct {
    Data []IpResponse   `json:"visits"`
    Meta PaginationMeta `json:"meta"`
}

type PaginatedIpHistoryResponse struct {
    Data []IpResponse   `json:"histories"`
    Meta PaginationMeta `json:"meta"`
}