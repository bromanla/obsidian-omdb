package omdb

type MovieResponse struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	Poster     string `json:"Poster"`
	Plot       string `json:"Plot"`
	ImdbRating string `json:"imdbRating"`
	Runtime    string `json:"Runtime"`
	Genre      string `json:"Genre"`
	Actors     string `json:"Actors"`
	Director   string `json:"Director"`
	Response   string `json:"Response"`
	Error      string `json:"Error"`
}
