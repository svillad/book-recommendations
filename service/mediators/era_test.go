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

type EraStoreMock struct {
	EraField   []models.Era
	ErrorField error
}

func (m *EraStoreMock) GetAllEras(ctx context.Context) ([]models.Era, error) {
	return m.EraField, m.ErrorField
}

func TestEraController_Get(t *testing.T) {
	MaxYear := int64(1969)
	MinYear := int64(1970)

	var cases = []struct {
		name   string
		store  stores.EraStore
		assert func(eras []models.Era, err error)
	}{
		{
			name: "success",
			store: &EraStoreMock{
				EraField: []models.Era{
					{
						ID:      1,
						Title:   "Classic",
						MaxYear: &MaxYear,
					},
					{
						ID:      2,
						Title:   "Modern",
						MinYear: &MinYear,
					},
				},
				ErrorField: nil,
			},
			assert: func(eras []models.Era, err error) {
				assert.Nil(t, err)
				assert.Len(t, eras, 2)
				assert.Equal(t, eras[0].ID, int64(1))
				assert.Equal(t, eras[0].Title, "Classic")
				assert.Nil(t, eras[0].MinYear)
				assert.Equal(t, eras[0].MaxYear, &MaxYear)
				assert.Equal(t, eras[1].ID, int64(2))
				assert.Equal(t, eras[1].Title, "Modern")
				assert.Equal(t, eras[1].MinYear, &MinYear)
				assert.Nil(t, eras[1].MaxYear)
			},
		},
		{
			name: "failure",
			store: &EraStoreMock{
				EraField:   []models.Era{},
				ErrorField: errors.New("Error"),
			},
			assert: func(eras []models.Era, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, c := range cases {
		m := mediators.NewEraMediator(log.NewEntry(log.New()), c.store)
		res, err := m.Get(context.Background())
		c.assert(res, err)
	}
}
