package translators

import (
	"net/http"

	"github.com/book-recommendations/service/models"
)

const (
	authorsParam  string = "authors"   //string
	genresParam   string = "genres"    //string
	minPagesParam string = "min-pages" //integer
	maxPagesParam string = "max-pages" //integer
	minYearParam  string = "min-year"  //integer
	maxYearParam  string = "max-year"  //integer
	limitParam    string = "limit"     //integer
)

// ToBookRequest creates the BookRequest model from the data in the request
func ToBooksRequest(r *http.Request) models.BookRequest {
	query := r.URL.Query()

	return models.BookRequest{
		Authors:  query.Get(authorsParam),
		Genres:   query.Get(genresParam),
		MinPages: query.Get(minPagesParam),
		MaxPages: query.Get(maxPagesParam),
		MinYear:  query.Get(minYearParam),
		MaxYear:  query.Get(maxYearParam),
		Limit:    query.Get(limitParam),
	}
}
