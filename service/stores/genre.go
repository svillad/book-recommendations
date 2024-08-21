package stores

import (
	"context"
	"fmt"

	"github.com/book-recommendations/service/models"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	tableGenre = "genre"
)

// GenreStore specifies the methods to get genres
type GenreStore interface {
	GetAllGenres(ctx context.Context) ([]models.Genre, error)
}

type genreStore struct {
	logger *log.Entry
	db     *sqlx.DB
}

func NewGenreStore(logger *log.Entry, db *sqlx.DB) GenreStore {
	return &genreStore{
		logger: logger,
		db:     db,
	}
}

func (s *genreStore) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	getGenresSQL := fmt.Sprintf(`SELECT id, title FROM %s`, tableGenre)

	rows, err := s.db.QueryContext(ctx, getGenresSQL)
	if err != nil {
		return nil, fmt.Errorf("error while building query: %w", err)
	}
	defer func() {
		errClose := rows.Close()
		errRows := rows.Err()
		if errClose != nil || errRows != nil {
			s.logger.WithFields(log.Fields{
				"errClose": errClose,
				"errRows":  errRows,
			}).Error("something went wrong while closing rows")
		}
	}()
	genres := make([]models.Genre, 0)
	for rows.Next() {
		var (
			id    int64
			title string
		)
		if err := rows.Scan(&id, &title); err != nil {
			return nil, fmt.Errorf("error getting genres: %w", err)
		}
		genres = append(
			genres,
			models.Genre{
				ID:    id,
				Title: title,
			},
		)
	}

	return genres, nil
}
