package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/book-recommendations/service/controllers/translators"
	"github.com/book-recommendations/service/mediators"
	log "github.com/sirupsen/logrus"
)

// EraController defines the controller for era data
type EraController struct {
	Logger             *log.Entry
	EraMediatorFactory func() mediators.EraMediator
}

// Get retrieves eras from the books backend
func (c *EraController) Get(w http.ResponseWriter, r *http.Request) {
	c.Logger.WithField("url", r.URL).Info("request")

	eraMediator := c.EraMediatorFactory()
	eras, err := eraMediator.Get(context.Background())
	if err != nil {
		c.Logger.WithError(err).Error("internal server error")
		translators.ParseError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(eras)
}
