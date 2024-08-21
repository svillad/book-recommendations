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

type AuthorStoreMock struct {
	AuthorField []models.Author
	ErrorField  error
}

func (m *AuthorStoreMock) GetAllAuthors(ctx context.Context) ([]models.Author, error) {
	return m.AuthorField, m.ErrorField
}

func TestAuthorController_Get(t *testing.T) {
	var cases = []struct {
		name   string
		store  stores.AuthorStore
		assert func(authors []models.Author, err error)
	}{
		{
			name: "success",
			store: &AuthorStoreMock{
				AuthorField: []models.Author{
					{
						ID:        1,
						FirstName: "Abraham",
						LastName:  "Stackhouse",
					},
					{
						ID:        2,
						FirstName: "Amelia",
						LastName:  "Wangerin, Jr.",
					},
					{
						ID:        3,
						FirstName: "Anastasia",
						LastName:  "Inez",
					},
				},
				ErrorField: nil,
			},
			assert: func(authors []models.Author, err error) {
				assert.Nil(t, err)
				assert.Len(t, authors, 3)
				assert.Equal(t, authors[0].ID, int64(1))
				assert.Equal(t, authors[0].FirstName, "Abraham")
				assert.Equal(t, authors[0].LastName, "Stackhouse")
			},
		},
		{
			name: "failure",
			store: &AuthorStoreMock{
				AuthorField: []models.Author{},
				ErrorField:  errors.New("Error"),
			},
			assert: func(authors []models.Author, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, c := range cases {
		m := mediators.NewAuthorMediator(log.NewEntry(log.New()), c.store)
		res, err := m.Get(context.Background())
		c.assert(res, err)
	}
}
