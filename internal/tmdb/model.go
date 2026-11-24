package tmdb

// Result is a movie or TV show returned by a search.
type Result struct {
	ID    int
	Title string
	Year  string
}

// Details holds the IMDb ID and localized text for a movie or TV show.
type Details struct {
	ImdbID string
	Title  string
	Plot   string
}
