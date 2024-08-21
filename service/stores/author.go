package stores

import (
	"context"
	"fmt"

	"github.com/book-recommendations/service/models"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const (
	tableAuthor = "author"
)

// AuthorStore specifies the methods to get authors
type AuthorStore interface {
	GetAllAuthors(ctx context.Context) ([]models.Author, error)
}

type authorStore struct {
	logger *log.Entry
	db     *sqlx.DB
}

func NewAuthorStore(logger *log.Entry, db *sqlx.DB) AuthorStore {
	return &authorStore{
		logger: logger,
		db:     db,
	}
}

func (s *authorStore) GetAllAuthors(ctx context.Context) ([]models.Author, error) {
	getAuthorsSQL := fmt.Sprintf(`SELECT id, first_name, last_name FROM %s`, tableAuthor)

	rows, err := s.db.QueryContext(ctx, getAuthorsSQL)
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
	authors := make([]models.Author, 0)
	for rows.Next() {
		var (
			id        int64
			firstName string
			lastName  string
		)
		if err := rows.Scan(&id, &firstName, &lastName); err != nil {
			return nil, fmt.Errorf("error getting authors: %w", err)
		}
		authors = append(
			authors,
			models.Author{
				ID:        id,
				FirstName: firstName,
				LastName:  lastName,
			},
		)
	}

	return authors, nil
}
