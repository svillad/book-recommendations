package mediators

import (
	"context"

	"github.com/book-recommendations/service/models"
	"github.com/book-recommendations/service/stores"
	log "github.com/sirupsen/logrus"
)

// AuthorMediator specifies the methods to get authors
type AuthorMediator interface {
	Get(ctx context.Context) ([]models.Author, error)
}

// authorMediator is the concrete implementation of the AuthorMediator interface
type authorMediator struct {
	logger *log.Entry
	store  stores.AuthorStore
}

// NewAuthorMediator returns a new instance of AuthorMediator
func NewAuthorMediator(logger *log.Entry, authorStore stores.AuthorStore) AuthorMediator {
	return &authorMediator{
		logger: logger,
		store:  authorStore,
	}
}

// Get returns a list of Authors
func (m *authorMediator) Get(ctx context.Context) ([]models.Author, error) {
	var authors []models.Author

	authors, err := m.store.GetAllAuthors(ctx)

	if err != nil {
		return nil, err
	}

	return authors, nil
}
