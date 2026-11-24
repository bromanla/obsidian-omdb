package tmdb

import (
	"cmp"
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var imdbIDPattern = regexp.MustCompile(`^tt[0-9]+$`)

type detailsResponse struct {
	OriginalLanguage string `json:"original_language"`
	OriginalTitle    string `json:"original_title"`
	OriginalName     string `json:"original_name"`
	ExternalIDs      struct {
		ImdbID string `json:"imdb_id"`
	} `json:"external_ids"`
	Translations struct {
		Items []struct {
			Language string `json:"iso_639_1"`
			Country  string `json:"iso_3166_1"`
			Data     struct {
				Title    string `json:"title"`
				Name     string `json:"name"`
				Overview string `json:"overview"`
			} `json:"data"`
		} `json:"translations"`
	} `json:"translations"`
}

func (c *Client) GetDetails(ctx context.Context, kind string, id int) (Details, error) {
	if kind != Movie && kind != TV {
		return Details{}, fmt.Errorf("unsupported TMDB type %q", kind)
	}

	params := url.Values{"append_to_response": {"translations,external_ids"}}
	var response detailsResponse
	if err := c.get(ctx, endpoint.JoinPath(kind, strconv.Itoa(id)), params, &response); err != nil {
		return Details{}, err
	}
	imdbID := response.ExternalIDs.ImdbID
	if !imdbIDPattern.MatchString(imdbID) {
		return Details{}, fmt.Errorf("TMDB %s %d has no valid IMDb ID", kind, id)
	}

	details := Details{ImdbID: imdbID}
	originalTitle := response.OriginalTitle
	if kind == TV {
		originalTitle = response.OriginalName
	}

	// TMDB can return fallback text in the top-level fields.
	language, country, _ := strings.Cut(c.language, "-")
	for _, translation := range response.Translations.Items {
		if !strings.EqualFold(translation.Language, language) || (country != "" && !strings.EqualFold(translation.Country, country)) {
			continue
		}

		title := translation.Data.Title
		if kind == TV {
			title = translation.Data.Name
		}
		details.Title = cmp.Or(details.Title, strings.TrimSpace(title))
		details.Plot = cmp.Or(details.Plot, strings.TrimSpace(translation.Data.Overview))
		if details.Title != "" && details.Plot != "" {
			break
		}
	}

	if details.Title == "" && strings.EqualFold(response.OriginalLanguage, language) {
		details.Title = strings.TrimSpace(originalTitle)
	}
	return details, nil
}
