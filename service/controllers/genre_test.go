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

type GenreMediatorMock struct {
	GenreField []models.Genre
	ErrorField error
}

func (m *GenreMediatorMock) Get(ctx context.Context) ([]models.Genre, error) {
	return m.GenreField, m.ErrorField
}

func TestGenreController_Get(t *testing.T) {
	var cases = []struct {
		name           string
		genreMediators mediators.GenreMediator
		request        interface{}
		assert         func(resp *http.Response, genres []models.Genre)
	}{
		{
			name: "success",
			genreMediators: &GenreMediatorMock{
				GenreField: []models.Genre{
					{
						ID:    1,
						Title: "Young Adult",
					},
					{
						ID:    2,
						Title: "SciFi/Fantasy",
					},
					{
						ID:    3,
						Title: "Romance",
					},
				},
				ErrorField: nil,
			},
			assert: func(resp *http.Response, genres []models.Genre) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Len(t, genres, 3)
				assert.Equal(t, genres[0].ID, int64(1))
				assert.Equal(t, genres[0].Title, "Young Adult")
			},
		},
		{
			name: "failure",
			genreMediators: &GenreMediatorMock{
				GenreField: []models.Genre{},
				ErrorField: errors.New("Error"),
			},
			assert: func(resp *http.Response, genres []models.Genre) {
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}
	for _, c := range cases {
		genreMediatorFactory := func() mediators.GenreMediator {
			return c.genreMediators
		}
		controller := controllers.GenreController{
			Logger:               log.NewEntry(log.New()),
			GenreMediatorFactory: genreMediatorFactory,
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "http://test.com/api/v1/genres", nil)

		router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
		router.HandleFunc("/genres", controller.Get).Methods(http.MethodGet)

		router.ServeHTTP(recorder, request)

		resp := recorder.Result()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "should return a readable response body")
		responseBody := []models.Genre{}
		err = json.Unmarshal(body, &responseBody)

		c.assert(resp, responseBody)
	}
}
