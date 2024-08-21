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

type GenreStoreMock struct {
	GenreField []models.Genre
	ErrorField error
}

func (m *GenreStoreMock) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	return m.GenreField, m.ErrorField
}

func TestGenreController_Get(t *testing.T) {
	var cases = []struct {
		name   string
		store  stores.GenreStore
		assert func(genres []models.Genre, err error)
	}{
		{
			name: "success",
			store: &GenreStoreMock{
				GenreField: []models.Genre{
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
			assert: func(genres []models.Genre, err error) {
				assert.Nil(t, err)
				assert.Len(t, genres, 3)
				assert.Equal(t, genres[0].ID, int64(1))
				assert.Equal(t, genres[0].Title, "Young Adult")
			},
		},
		{
			name: "failure",
			store: &GenreStoreMock{
				GenreField: []models.Genre{},
				ErrorField: errors.New("Error"),
			},
			assert: func(genres []models.Genre, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, c := range cases {
		m := mediators.NewGenreMediator(log.NewEntry(log.New()), c.store)
		res, err := m.Get(context.Background())
		c.assert(res, err)
	}
}
