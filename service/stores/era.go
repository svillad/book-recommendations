package stores

import (
	"context"
	"fmt"

	"github.com/book-recommendations/service/models"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	tableEra = "era"
)

// EraStore specifies the methods to get eras
type EraStore interface {
	GetAllEras(ctx context.Context) ([]models.Era, error)
}

type eraStore struct {
	logger *log.Entry
	db     *sqlx.DB
}

func NewEraStore(logger *log.Entry, db *sqlx.DB) EraStore {
	return &eraStore{
		logger: logger,
		db:     db,
	}
}

func (s *eraStore) GetAllEras(ctx context.Context) ([]models.Era, error) {
	getErasSQL := fmt.Sprintf(`SELECT id, title, min_year, max_year FROM %s`, tableEra)

	rows, err := s.db.QueryContext(ctx, getErasSQL)
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
	eras := make([]models.Era, 0)
	for rows.Next() {
		var (
			id      int64
			title   string
			minYear *int64
			maxYear *int64
		)
		if err := rows.Scan(&id, &title, &minYear, &maxYear); err != nil {
			return nil, fmt.Errorf("error getting eras: %w", err)
		}

		eras = append(
			eras,
			models.Era{
				ID:      id,
				Title:   title,
				MinYear: minYear,
				MaxYear: maxYear,
			},
		)
	}

	return eras, nil
}
