// Package catalog searches TMDB and combines its translations with OMDb metadata.
package catalog

import (
	"cmp"
	"context"

	"obsidian/omdb/internal/omdb"
	"obsidian/omdb/internal/tmdb"
)

type Service struct {
	tmdb *tmdb.Client
	omdb *omdb.Client
}

func New(tmdbClient *tmdb.Client, omdbClient *omdb.Client) *Service {
	return &Service{tmdb: tmdbClient, omdb: omdbClient}
}

func (s *Service) Search(ctx context.Context, query Query, kind Kind) ([]tmdb.Result, error) {
	return s.tmdb.Search(ctx, string(kind), query.Title, query.Year)
}

// Details builds a card using TMDB translations and OMDb metadata.
func (s *Service) Details(ctx context.Context, kind Kind, tmdbID int) (Card, error) {
	localized, err := s.tmdb.GetDetails(ctx, string(kind), tmdbID)
	if err != nil {
		return Card{}, err
	}

	metadata, err := s.omdb.GetByID(ctx, localized.ImdbID)
	if err != nil {
		return Card{}, err
	}
	return Card{
		Title:      cmp.Or(localized.Title, metadata.Title),
		Year:       metadata.Year,
		ImdbID:     localized.ImdbID,
		Kind:       kind,
		Poster:     metadata.Poster,
		Plot:       cmp.Or(localized.Plot, metadata.Plot),
		ImdbRating: metadata.ImdbRating,
		Runtime:    metadata.Runtime,
		Genre:      metadata.Genre,
		Actors:     metadata.Actors,
		Director:   metadata.Director,
	}, nil
}
