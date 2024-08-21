package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/book-recommendations/service/controllers/translators"
	"github.com/book-recommendations/service/mediators"
	log "github.com/sirupsen/logrus"
)

// GenreController defines the controller for genre data
type GenreController struct {
	Logger               *log.Entry
	GenreMediatorFactory func() mediators.GenreMediator
}

// Get retrieves genre from the books backend
func (c *GenreController) Get(w http.ResponseWriter, r *http.Request) {
	c.Logger.WithField("url", r.URL).Info("request")

	genreMediator := c.GenreMediatorFactory()
	genres, err := genreMediator.Get(context.Background())
	if err != nil {
		c.Logger.WithError(err).Error("internal server error")
		translators.ParseError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}
