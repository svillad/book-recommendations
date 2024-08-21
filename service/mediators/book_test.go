package mediators_test

import (
	"context"
	"errors"
	"testing"

	"github.com/book-recommendations/service/mediators"
	"github.com/book-recommendations/service/models"
	"github.com/book-recommendations/service/stores"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type BookStoreMock struct {
	BookField  []models.Book
	ErrorField error
}

func (m *BookStoreMock) GetBooks(ctx context.Context, req models.BookRequest) ([]models.Book, error) {
	return m.BookField, m.ErrorField
}

func TestBookController_Get(t *testing.T) {
	var cases = []struct {
		name    string
		store   stores.BookStore
		request models.BookRequest
		assert  func(books []models.Book, err error)
	}{
		{
			name: "success",
			store: &BookStoreMock{
				BookField: []models.Book{
					{
						ID:    1,
						Title: "Young Adult",
					},
					{
						ID:    2,
						Title: "SciFi/Fantasy",
					},
					{
						ID:    3,
						Title: "Romance",
					},
				},
				ErrorField: nil,
			},
			request: models.BookRequest{},
			assert: func(books []models.Book, err error) {
				assert.Nil(t, err)
				assert.Len(t, books, 3)
				assert.Equal(t, books[0].ID, int64(1))
				assert.Equal(t, books[0].Title, "Young Adult")
			},
		},
		{
			name: "failure",
			store: &BookStoreMock{
				BookField:  []models.Book{},
				ErrorField: errors.New("Error"),
			},
			request: models.BookRequest{},
			assert: func(books []models.Book, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, c := range cases {
		m := mediators.NewBookMediator(log.NewEntry(log.New()), c.store)
		res, err := m.Get(context.Background(), c.request)
		c.assert(res, err)
	}
}
