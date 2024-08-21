package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/book-recommendations/service/controllers/translators"
	"github.com/book-recommendations/service/mediators"
	log "github.com/sirupsen/logrus"
)

// AuthorController defines the controller for author data
type AuthorController struct {
	Logger                *log.Entry
	AuthorMediatorFactory func() mediators.AuthorMediator
}

// Get retrieves authors from the books backend
func (c *AuthorController) Get(w http.ResponseWriter, r *http.Request) {
	c.Logger.WithField("url", r.URL).Info("request")

	authorMediator := c.AuthorMediatorFactory()
	authors, err := authorMediator.Get(context.Background())
	if err != nil {
		c.Logger.WithError(err).Error("internal server error")
		translators.ParseError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}
