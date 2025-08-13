package omdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	OmdbUrl = "http://www.omdbapi.com/"
)

var imdbIdPattern = regexp.MustCompile(`^tt\d+$`)

type Client struct {
	apiKey string
	client *http.Client
}

func New(apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{
		apiKey: apiKey,
		client: httpClient,
	}
}

func (c *Client) buildURL(params map[string]string) string {
	u, _ := url.Parse(OmdbUrl)
	q := u.Query()

	q.Add("apikey", c.apiKey)
	q.Add("r", "json")
	for key, value := range params {
		q.Add(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func (c *Client) get(ctx context.Context, params map[string]string, out any) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildURL(params), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (c *Client) IsImdbId(s string) bool {
	return imdbIdPattern.MatchString(s)
}

func (c *Client) GetByID(ctx context.Context, id string) (MovieResponse, error) {
	params := map[string]string{"i": id}
	var resp MovieResponse

	if err := c.get(ctx, params, &resp); err != nil {
		return resp, err
	}

	if resp.Response != "True" {
		return resp, fmt.Errorf("OMDb API error: %s", resp.Error)
	}

	return resp, nil
}

func (c *Client) GetByName(ctx context.Context, name string) (SearchResponse, error) {
	params := map[string]string{"s": name}
	var resp SearchResponse

	if err := c.get(ctx, params, &resp); err != nil {
		return resp, err
	}

	if resp.Response != "True" {
		return resp, fmt.Errorf("OMDb API error: %s", resp.Error)
	}

	return resp, nil
}

func (c *Client) Find(ctx context.Context, query string) ([]MovieShort, error) {

	if c.IsImdbId(query) {
		movie, err := c.GetByID(ctx, query)
		if err != nil {
			return nil, err
		}

		return []MovieShort{{
			Title:  movie.Title,
			Year:   movie.Year,
			ImdbID: movie.ImdbID,
			Type:   movie.Type,
			Poster: movie.Poster,
		}}, nil
	}

	movies, err := c.GetByName(ctx, query)
	if err != nil {
		return nil, err
	}

	return movies.Search, nil
}
