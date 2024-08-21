package mediators

import (
	"context"

	"github.com/book-recommendations/service/models"
	"github.com/book-recommendations/service/stores"
	log "github.com/sirupsen/logrus"
)

// SizeMediator specifies the methods to get sizes
type SizeMediator interface {
	Get(ctx context.Context) ([]models.Size, error)
}

// sizeMediator is the concrete implementation of the SizeMediator interface
type sizeMediator struct {
	logger *log.Entry
	store  stores.SizeStore
}

// NewSizeMediator returns a new instance of SizeMediator
func NewSizeMediator(logger *log.Entry, sizeStore stores.SizeStore) SizeMediator {
	return &sizeMediator{
		logger: logger,
		store:  sizeStore,
	}
}

// Get returns a list of Sizes
func (m *sizeMediator) Get(ctx context.Context) ([]models.Size, error) {
	var sizes []models.Size

	sizes, err := m.store.GetAllSizes(ctx)
	if err != nil {
		return nil, err
	}

	return sizes, nil
}
