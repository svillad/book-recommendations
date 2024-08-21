package mediators

import (
	"context"

	"github.com/book-recommendations/service/models"
	"github.com/book-recommendations/service/stores"
	log "github.com/sirupsen/logrus"
)

// BookMediator specifies the methods to get books
type BookMediator interface {
	Get(ctx context.Context, req models.BookRequest) ([]models.Book, error)
}

// bookMediator is the concrete implementation of the BookMediator interface
type bookMediator struct {
	logger *log.Entry
	store  stores.BookStore
}

// NewBookMediator returns a new instance of BookMediator
func NewBookMediator(logger *log.Entry, bookStore stores.BookStore) BookMediator {
	return &bookMediator{
		logger: logger,
		store:  bookStore,
	}
}

// Get returns a list of Books
func (m *bookMediator) Get(ctx context.Context, req models.BookRequest) ([]models.Book, error) {
	var books []models.Book

	books, err := m.store.GetBooks(ctx, req)
	if err != nil {
		return nil, err
	}

	return books, nil
}
