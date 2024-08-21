package mediators

import (
	"context"

	"github.com/book-recommendations/service/models"
	"github.com/book-recommendations/service/stores"
	log "github.com/sirupsen/logrus"
)

// EraMediator specifies the methods to get eras
type EraMediator interface {
	Get(ctx context.Context) ([]models.Era, error)
}

// eraMediator is the concrete implementation of the EraMediator interface
type eraMediator struct {
	logger *log.Entry
	store  stores.EraStore
}

// NewEraMediator returns a new instance of EraMediator
func NewEraMediator(logger *log.Entry, eraStore stores.EraStore) EraMediator {
	return &eraMediator{
		logger: logger,
		store:  eraStore,
	}
}

// Get returns a list of Eras
func (m *eraMediator) Get(ctx context.Context) ([]models.Era, error) {
	var eras []models.Era

	eras, err := m.store.GetAllEras(ctx)
	if err != nil {
		return nil, err
	}

	return eras, nil
}
