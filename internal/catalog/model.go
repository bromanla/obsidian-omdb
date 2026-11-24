package catalog

import "obsidian/omdb/internal/tmdb"

// Kind identifies a movie or a TV show.
type Kind string

const (
	Movie  Kind = tmdb.Movie
	TVShow Kind = tmdb.TV
)

func (k Kind) Label() string {
	if k == TVShow {
		return "TV show"
	}
	return "Movie"
}

// Query holds the title and optional year to search for.
type Query struct {
	Title string
	Year  string
}

// Card contains the fields used by the Telegram preview and Markdown template.
type Card struct {
	Title      string
	Year       string
	ImdbID     string
	Kind       Kind
	Poster     string
	Plot       string
	ImdbRating string
	Runtime    string
	Genre      string
	Actors     string
	Director   string
}
