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

type EraMediatorMock struct {
	EraField   []models.Era
	ErrorField error
}

func (m *EraMediatorMock) Get(ctx context.Context) ([]models.Era, error) {
	return m.EraField, m.ErrorField
}

func TestEraController_Get(t *testing.T) {
	MaxYear := int64(1969)
	MinYear := int64(1970)

	var cases = []struct {
		name         string
		eraMediators mediators.EraMediator
		request      interface{}
		assert       func(resp *http.Response, eras []models.Era)
	}{
		{
			name: "success",
			eraMediators: &EraMediatorMock{
				EraField: []models.Era{
					{
						ID:      1,
						Title:   "Classic",
						MaxYear: &MaxYear,
					},
					{
						ID:      2,
						Title:   "Modern",
						MinYear: &MinYear,
					},
				},
				ErrorField: nil,
			},
			assert: func(resp *http.Response, eras []models.Era) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Len(t, eras, 2)
				assert.Equal(t, eras[0].ID, int64(1))
				assert.Equal(t, eras[0].Title, "Classic")
				assert.Nil(t, eras[0].MinYear)
				assert.Equal(t, eras[0].MaxYear, &MaxYear)
				assert.Equal(t, eras[1].ID, int64(2))
				assert.Equal(t, eras[1].Title, "Modern")
				assert.Equal(t, eras[1].MinYear, &MinYear)
				assert.Nil(t, eras[1].MaxYear)
			},
		},
		{
			name: "failure",
			eraMediators: &EraMediatorMock{
				EraField:   []models.Era{},
				ErrorField: errors.New("Error"),
			},
			assert: func(resp *http.Response, eras []models.Era) {
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}
	for _, c := range cases {
		eraMediatorFactory := func() mediators.EraMediator {
			return c.eraMediators
		}
		controller := controllers.EraController{
			Logger:             log.NewEntry(log.New()),
			EraMediatorFactory: eraMediatorFactory,
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "http://test.com/api/v1/eras", nil)

		router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
		router.HandleFunc("/eras", controller.Get).Methods(http.MethodGet)

		router.ServeHTTP(recorder, request)

		resp := recorder.Result()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "should return a readable response body")
		responseBody := []models.Era{}
		err = json.Unmarshal(body, &responseBody)

		c.assert(resp, responseBody)
	}
}
