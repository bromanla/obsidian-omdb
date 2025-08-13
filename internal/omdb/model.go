package omdb

import (
	"fmt"
	"regexp"
	"strings"
)

type ApiResponse struct {
	Response string `json:"Response"`
	Error    string `json:"Error"`
}

type MovieShort struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	ImdbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}

// Header returns a formatted string for display, e.g., "(M) Title (Year)".
func (m MovieShort) Header() string {
	tag := strings.ToUpper(m.Type[0:1])
	timeline := m.Year
	if strings.HasSuffix(timeline, "–") {
		timeline += "..."
	}

	return fmt.Sprintf("(%s) %s (%s)", tag, m.Title, timeline)
}

type Rating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

type SearchResponse struct {
	ApiResponse
	Search       []MovieShort `json:"Search"`
	TotalResults string       `json:"totalResults"`
}

type MovieResponse struct {
	Rated      string   `json:"Rated"`
	Released   string   `json:"Released"`
	Runtime    string   `json:"Runtime"`
	Genre      string   `json:"Genre"`
	Director   string   `json:"Director"`
	Writer     string   `json:"Writer"`
	Actors     string   `json:"Actors"`
	Plot       string   `json:"Plot"`
	Language   string   `json:"Language"`
	Country    string   `json:"Country"`
	Awards     string   `json:"Awards"`
	Ratings    []Rating `json:"Ratings"`
	Metascore  string   `json:"Metascore"`
	ImdbRating string   `json:"imdbRating"`
	ImdbVotes  string   `json:"imdbVotes"`
	DVD        string   `json:"DVD"`
	BoxOffice  string   `json:"BoxOffice"`
	Production string   `json:"Production"`
	Website    string   `json:"Website"`

	MovieShort
	ApiResponse
}

// Define invalid characters pattern
var invalidChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)

// Sanitize returns a sanitized version of the movie title for use in filenames.
func (m MovieResponse) Sanitize() string {
	filename := m.Title
	// Replace invalid characters with an underscore
	filename = invalidChars.ReplaceAllString(filename, "")
	// Remove leading/trailing spaces and dots
	filename = strings.TrimSpace(filename)

	return filename
}
