package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var endpoint = url.URL{
	Scheme: "https",
	Host:   "api.themoviedb.org",
	Path:   "/3",
}

const (
	Movie = "movie"
	TV    = "tv"
)

type Client struct {
	apiKey   string
	language string
	http     *http.Client
}

func New(apiKey, language string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{apiKey: apiKey, language: language, http: httpClient}
}

func (c *Client) get(ctx context.Context, requestURL *url.URL, params url.Values, out any) error {
	params.Set("language", c.language)
	requestURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return fmt.Errorf("create TMDB request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("send TMDB request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TMDB status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode TMDB response: %w", err)
	}
	return nil
}
