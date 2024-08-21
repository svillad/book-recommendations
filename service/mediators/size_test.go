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

type SizeStoreMock struct {
	SizeField  []models.Size
	ErrorField error
}

func (m *SizeStoreMock) GetAllSizes(ctx context.Context) ([]models.Size, error) {
	return m.SizeField, m.ErrorField
}

func TestSizeController_Get(t *testing.T) {
	MaxPages1 := int64(34)
	MinPages1 := int64(35)
	MaxPages2 := int64(84)
	MinPages2 := int64(800)

	var cases = []struct {
		name   string
		store  stores.SizeStore
		assert func(sizes []models.Size, err error)
	}{
		{
			name: "success",
			store: &SizeStoreMock{
				SizeField: []models.Size{
					{
						ID:       1,
						Title:    "Short story – up to 35 pages",
						MaxPages: &MaxPages1,
					},
					{
						ID:       2,
						Title:    "Novelette – 35 to 85 pages",
						MinPages: &MinPages1,
						MaxPages: &MaxPages2,
					},
					{
						ID:       3,
						Title:    "Monument – 800 pages and up",
						MinPages: &MinPages2,
					},
				},
				ErrorField: nil,
			},
			assert: func(sizes []models.Size, err error) {
				assert.Nil(t, err)
				assert.Len(t, sizes, 3)
				assert.Equal(t, sizes[0].ID, int64(1))
				assert.Equal(t, sizes[0].Title, "Short story – up to 35 pages")
				assert.Nil(t, sizes[0].MinPages)
				assert.Equal(t, sizes[0].MaxPages, &MaxPages1)
			},
		},
		{
			name: "failure",
			store: &SizeStoreMock{
				SizeField:  []models.Size{},
				ErrorField: errors.New("Error"),
			},
			assert: func(sizes []models.Size, err error) {
				assert.NotNil(t, err)
			},
		},
	}
	for _, c := range cases {
		m := mediators.NewSizeMediator(log.NewEntry(log.New()), c.store)
		res, err := m.Get(context.Background())
		c.assert(res, err)
	}
}
