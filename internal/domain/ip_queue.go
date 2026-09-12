package domain

type IpQueueEntity struct {
    ID        string `json:"id"`
    IpAddress string `json:"ip_address"`
    UserAgent string `json:"user_agent"`
    Attempt   int    `json:"attempt"`
}