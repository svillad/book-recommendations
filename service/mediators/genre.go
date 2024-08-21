package mediators

import (
	"context"

	"github.com/book-recommendations/service/models"
	"github.com/book-recommendations/service/stores"
	log "github.com/sirupsen/logrus"
)

// GenreMediator specifies the methods to get genres
type GenreMediator interface {
	Get(ctx context.Context) ([]models.Genre, error)
}

// genreMediator is the concrete implementation of the GenreMediator interface
type genreMediator struct {
	logger *log.Entry
	store  stores.GenreStore
}

// NewGenreMediator returns a new instance of GenreMediator
func NewGenreMediator(logger *log.Entry, genreStore stores.GenreStore) GenreMediator {
	return &genreMediator{
		logger: logger,
		store:  genreStore,
	}
}

// Get returns a list of Genres
func (m *genreMediator) Get(ctx context.Context) ([]models.Genre, error) {
	var genres []models.Genre

	genres, err := m.store.GetAllGenres(ctx)
	if err != nil {
		return nil, err
	}

	return genres, nil
}
