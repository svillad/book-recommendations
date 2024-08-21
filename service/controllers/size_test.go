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

type SizeMediatorMock struct {
	SizeField  []models.Size
	ErrorField error
}

func (m *SizeMediatorMock) Get(ctx context.Context) ([]models.Size, error) {
	return m.SizeField, m.ErrorField
}

func TestSizeController_Get(t *testing.T) {
	MaxPages1 := int64(34)
	MinPages1 := int64(35)
	MaxPages2 := int64(84)
	MinPages2 := int64(800)

	var cases = []struct {
		name          string
		sizeMediators mediators.SizeMediator
		request       interface{}
		assert        func(resp *http.Response, sizes []models.Size)
	}{
		{
			name: "success",
			sizeMediators: &SizeMediatorMock{
				SizeField: []models.Size{
					{
						ID:       1,
						Title:    "Short story – up to 35 pages",
						MaxPages: &MaxPages1,
					},
					{
						ID:       2,
						Title:    "Novelette – 35 to 85 pages",
						MinPages: &MinPages1,
						MaxPages: &MaxPages2,
					},
					{
						ID:       3,
						Title:    "Monument – 800 pages and up",
						MinPages: &MinPages2,
					},
				},
				ErrorField: nil,
			},
			assert: func(resp *http.Response, sizes []models.Size) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				assert.Len(t, sizes, 3)
				assert.Equal(t, sizes[0].ID, int64(1))
				assert.Equal(t, sizes[0].Title, "Short story – up to 35 pages")
				assert.Nil(t, sizes[0].MinPages)
				assert.Equal(t, sizes[0].MaxPages, &MaxPages1)
			},
		},
		{
			name: "failure",
			sizeMediators: &SizeMediatorMock{
				SizeField:  []models.Size{},
				ErrorField: errors.New("Error"),
			},
			assert: func(resp *http.Response, sizes []models.Size) {
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
			},
		},
	}
	for _, c := range cases {
		sizeMediatorFactory := func() mediators.SizeMediator {
			return c.sizeMediators
		}
		controller := controllers.SizeController{
			Logger:              log.NewEntry(log.New()),
			SizeMediatorFactory: sizeMediatorFactory,
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "http://test.com/api/v1/sizes", nil)

		router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
		router.HandleFunc("/sizes", controller.Get).Methods(http.MethodGet)

		router.ServeHTTP(recorder, request)

		resp := recorder.Result()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err, "should return a readable response body")
		responseBody := []models.Size{}
		err = json.Unmarshal(body, &responseBody)

		c.assert(resp, responseBody)
	}
}
