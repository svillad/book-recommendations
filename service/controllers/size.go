package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/book-recommendations/service/controllers/translators"
	"github.com/book-recommendations/service/mediators"
	log "github.com/sirupsen/logrus"
)

// SizeController defines the controller for size data
type SizeController struct {
	Logger              *log.Entry
	SizeMediatorFactory func() mediators.SizeMediator
}

// Get retrieves sizes from the books backend
func (c *SizeController) Get(w http.ResponseWriter, r *http.Request) {
	c.Logger.WithField("url", r.URL).Info("request")

	sizeMediator := c.SizeMediatorFactory()
	sizes, err := sizeMediator.Get(context.Background())
	if err != nil {
		c.Logger.WithError(err).Error("internal server error")
		translators.ParseError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sizes)
}
