package omdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var endpoint = url.URL{
	Scheme: "https",
	Host:   "www.omdbapi.com",
	Path:   "/",
}

type Client struct {
	apiKey string
	http   *http.Client
}

func New(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{apiKey: apiKey, http: httpClient}
}

func (c *Client) GetByID(ctx context.Context, id string) (MovieResponse, error) {
	var movie MovieResponse

	params := url.Values{
		"apikey": {c.apiKey},
		"i":      {id},
		"r":      {"json"},
	}

	requestURL := endpoint
	requestURL.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return movie, fmt.Errorf("create OMDb request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		// The URL in a transport error contains the API key.
		if requestErr, ok := errors.AsType[*url.Error](err); ok {
			err = requestErr.Err
		}
		return movie, fmt.Errorf("send OMDb request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movie, fmt.Errorf("OMDb status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return movie, fmt.Errorf("decode OMDb response: %w", err)
	}

	if movie.Response != "True" {
		return movie, fmt.Errorf("OMDb API error: %s", movie.Error)
	}

	return movie, nil
}
