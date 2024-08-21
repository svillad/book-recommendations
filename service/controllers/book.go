package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/book-recommendations/service/controllers/translators"
	"github.com/book-recommendations/service/mediators"
	log "github.com/sirupsen/logrus"
)

// BookController defines the controller for books data
type BookController struct {
	Logger              *log.Entry
	BookMediatorFactory func() mediators.BookMediator
}

// Get retrieves books from the books backend
func (c *BookController) Get(w http.ResponseWriter, r *http.Request) {
	c.Logger.WithField("url", r.URL).Info("request")

	req := translators.ToBooksRequest(r)
	if err := req.Validate(); err != nil {
		c.Logger.WithField("authors", req.Authors).WithField("genres", req.Genres).
			WithField("min-pages", req.MinPages).WithField("max-pages", req.MaxPages).
			WithField("min-year", req.MinYear).WithField("max-year", req.MaxYear).
			WithField("limit", req.Limit).WithError(err).Error("invalid request params for get books")
		translators.ParseError(w, http.StatusBadRequest)
		return
	}

	bookMediator := c.BookMediatorFactory()
	books, err := bookMediator.Get(context.Background(), req)
	if err != nil {
		c.Logger.WithError(err).Error("internal server error ")
		translators.ParseError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
