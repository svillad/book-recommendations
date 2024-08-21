package stores

import (
	"context"
	"fmt"

	"github.com/book-recommendations/service/models"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	tableSize = "size"
)

// SizeStore specifies the methods to get sizes
type SizeStore interface {
	GetAllSizes(ctx context.Context) ([]models.Size, error)
}

type sizeStore struct {
	logger *log.Entry
	db     *sqlx.DB
}

func NewSizeStore(logger *log.Entry, db *sqlx.DB) SizeStore {
	return &sizeStore{
		logger: logger,
		db:     db,
	}
}

func (s *sizeStore) GetAllSizes(ctx context.Context) ([]models.Size, error) {
	getSizesSQL := fmt.Sprintf(`SELECT id, title, min_pages, max_pages FROM %s`, tableSize)

	rows, err := s.db.QueryContext(ctx, getSizesSQL)
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
	sizes := make([]models.Size, 0)
	for rows.Next() {
		var (
			id       int64
			title    string
			minPages *int64
			maxPages *int64
		)
		if err := rows.Scan(&id, &title, &minPages, &maxPages); err != nil {
			return nil, fmt.Errorf("error getting sizes: %w", err)
		}

		sizes = append(
			sizes,
			models.Size{
				ID:       id,
				Title:    title,
				MinPages: minPages,
				MaxPages: maxPages,
			},
		)
	}

	return sizes, nil
}
