package controllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/book-recommendations/service/controllers"
	"github.com/book-recommendations/service/mediators"
	"github.com/book-recommendations/service/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BookMediatorMock struct {
	BookField  []models.Book
	ErrorField error
}

func (m *BookMediatorMock) Get(ctx context.Context, req models.BookRequest) ([]models.Book, error) {
	return m.BookField, m.ErrorField
}

func TestBookController_Get(t *testing.T) {
	var cases = []struct {
		name          string
		bookMediators mediators.BookMediator
		request       string
		assert        func(resp *http.Response, books []models.Book)
	}{
		{
			name: "success",
			bookMediators: &BookMediatorMock{
				BookField: []models.Book{
					{
						ID:            1,
						Title:         "Alanna Saves the Day",
						YearPublished: 1972,
						Rating:        1.62,
						Pages:         169,
						Genre: models.Genre{
							ID:    8,
							Title: "Childrens",
						},
						Author: models.Author{
							ID:        6,
							FirstName: "Bernard",
							LastName:  "Hopf",
						},
					},
					{
						ID:            2,
						Title:         "Adventures of Kaya",
						YearPublished: 1999,
						Rating:        2.13,
						Pages:         619,
						Genre: models.Genre{
							ID:    8,
							Title: "Young Adult",
						},
						Author: models.Author{
							ID:        40,
							FirstName: "Ward",
							LastName:  "Haigh",
						},
					},
				},
				ErrorField: nil,
			},
			request: "limit=10",
			assert: func(resp *http.Response, books []models.Book) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Len(t, books, 2)
				assert.Equal(t, books[0].ID, int64(1))
				assert.Equal(t, books[0].Title, "Alanna Saves the Day")
				assert.Equal(t, books[0].YearPublished, int64(1972))
				assert.Equal(t, books[0].Rating, 1.62)
				assert.Equal(t, books[0].Pages, int64(169))
				assert.Equal(t, books[0].Genre.ID, int64(8))
				assert.Equal(t, books[0].Genre.Title, "Childrens")
				assert.Equal(t, books[0].Author.ID, int64(6))
				assert.Equal(t, books[0].Author.FirstName, "Bernard")
				assert.Equal(t, books[0].Author.LastName, "Hopf")
			},
		},
		{
			name: "failure",
			bookMediators: &BookMediatorMock{
				BookField:  []models.Book{},
				ErrorField: errors.New("Error"),
			},
			assert: func(resp *http.Response, books []models.Book) {
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
		{
			name: "bad request",
			bookMediators: &BookMediatorMock{
				BookField:  []models.Book{},
				ErrorField: errors.New("Error"),
			},
			request: "limit=0",
			assert: func(resp *http.Response, books []models.Book) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			},
		},
	}
	for _, c := range cases {
		bookMediatorFactory := func() mediators.BookMediator {
			return c.bookMediators
		}
		controller := controllers.BookController{
			Logger:              log.NewEntry(log.New()),
			BookMediatorFactory: bookMediatorFactory,
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "http://test.com/api/v1/books?"+c.request, nil)

		router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
		router.HandleFunc("/books", controller.Get).Methods(http.MethodGet)

		router.ServeHTTP(recorder, request)

		resp := recorder.Result()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "should return a readable response body")
		responseBody := []models.Book{}
		err = json.Unmarshal(body, &responseBody)

		c.assert(resp, responseBody)
	}
}
