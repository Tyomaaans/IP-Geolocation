package ipgeo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"ip-geo/internal/domain"
	"net/http"
	"net/url"
	"time"
)

type ipGeoClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewIpGeoClient(apiKey string) domain.IpGeoClient {
	return &ipGeoClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *ipGeoClient) FetchByIP(ctx context.Context, ip string) (*domain.IpGeoEntity, error) {
	params := url.Values{}
	params.Set("apiKey", c.apiKey)
	params.Set("ip", ip)

	fullURL := "https://api.ipgeolocation.io/v3/ipgeo?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ipgeo api error (%d): %s", resp.StatusCode, string(body))
	}

	var result domain.IpGeoEntity
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}