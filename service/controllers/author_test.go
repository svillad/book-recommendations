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

type AuthorMediatorMock struct {
	AuthorField []models.Author
	ErrorField  error
}

func (m *AuthorMediatorMock) Get(ctx context.Context) ([]models.Author, error) {
	return m.AuthorField, m.ErrorField
}

func TestAuthorController_Get(t *testing.T) {
	var cases = []struct {
		name            string
		authorMediators mediators.AuthorMediator
		request         interface{}
		assert          func(resp *http.Response, authors []models.Author)
	}{
		{
			name: "success",
			authorMediators: &AuthorMediatorMock{
				AuthorField: []models.Author{
					{
						ID:        1,
						FirstName: "Abraham",
						LastName:  "Stackhouse",
					},
					{
						ID:        2,
						FirstName: "Amelia",
						LastName:  "Wangerin, Jr.",
					},
					{
						ID:        3,
						FirstName: "Anastasia",
						LastName:  "Inez",
					},
				},
				ErrorField: nil,
			},
			assert: func(resp *http.Response, authors []models.Author) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Len(t, authors, 3)
				assert.Equal(t, authors[0].ID, int64(1))
				assert.Equal(t, authors[0].FirstName, "Abraham")
				assert.Equal(t, authors[0].LastName, "Stackhouse")
			},
		},
		{
			name: "failure",
			authorMediators: &AuthorMediatorMock{
				AuthorField: []models.Author{},
				ErrorField:  errors.New("Error"),
			},
			assert: func(resp *http.Response, authors []models.Author) {
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}
	for _, c := range cases {
		authorMediatorFactory := func() mediators.AuthorMediator {
			return c.authorMediators
		}
		controller := controllers.AuthorController{
			Logger:                log.NewEntry(log.New()),
			AuthorMediatorFactory: authorMediatorFactory,
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "http://test.com/api/v1/authors", nil)

		router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
		router.HandleFunc("/authors", controller.Get).Methods(http.MethodGet)

		router.ServeHTTP(recorder, request)

		resp := recorder.Result()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "should return a readable response body")
		responseBody := []models.Author{}
		err = json.Unmarshal(body, &responseBody)

		c.assert(resp, responseBody)
	}
}
