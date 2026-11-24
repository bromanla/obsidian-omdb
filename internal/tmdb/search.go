package tmdb

import (
	"cmp"
	"context"
	"fmt"
	"net/url"
	"strings"
)

type searchResponse struct {
	Results []struct {
		ID            int    `json:"id"`
		Title         string `json:"title"`
		OriginalTitle string `json:"original_title"`
		ReleaseDate   string `json:"release_date"`
		Name          string `json:"name"`
		OriginalName  string `json:"original_name"`
		FirstAirDate  string `json:"first_air_date"`
	} `json:"results"`
}

func (c *Client) Search(ctx context.Context, kind, title, year string) ([]Result, error) {
	var yearParam string

	switch kind {
	case Movie:
		yearParam = "primary_release_year"
	case TV:
		yearParam = "first_air_date_year"
	default:
		return nil, fmt.Errorf("unsupported TMDB type %q", kind)
	}

	params := url.Values{"query": {title}}
	if year != "" {
		params.Set(yearParam, year)
	}

	var response searchResponse
	if err := c.get(ctx, endpoint.JoinPath("search", kind), params, &response); err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(response.Results))
	for _, item := range response.Results {
		title, originalTitle, date := item.Title, item.OriginalTitle, item.ReleaseDate
		if kind == TV {
			title, originalTitle, date = item.Name, item.OriginalName, item.FirstAirDate
		}
		title = cmp.Or(strings.TrimSpace(title), strings.TrimSpace(originalTitle))
		if item.ID <= 0 || title == "" {
			continue
		}
		year, _, _ := strings.Cut(date, "-")
		results = append(results, Result{ID: item.ID, Title: title, Year: year})
	}
	return results, nil
}
