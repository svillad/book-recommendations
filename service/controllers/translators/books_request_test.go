package translators_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/book-recommendations/service/controllers/translators"
	"github.com/book-recommendations/service/models"
	"github.com/stretchr/testify/assert"
)

func TestToBooksRequest(t *testing.T) {
	cases := []struct {
		name   string
		url    string
		assert func(resp models.BookRequest, err error)
	}{
		{
			name: "success",
			url:  "authors=1,11&genres=1&min-pages=100&max-pages=300&min-year=2001&max-year=2100&limit=5",
			assert: func(resp models.BookRequest, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "success - omit fields",
			url:  "",
			assert: func(resp models.BookRequest, err error) {
				assert.Nil(t, err)
			},
		},
		{
			name: "failure - invalid authors",
			url:  "authors=1,invalid&genres=1&min-pages=100&max-pages=300&min-year=2001&max-year=2100&limit=5",
			assert: func(resp models.BookRequest, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "authors: should be for exaple: 123,456,789.")
			},
		},
		{
			name: "failure - invalid genres",
			url:  "authors=1,11&genres=invalid&min-pages=100&max-pages=300&min-year=2001&max-year=2100&limit=5",
			assert: func(resp models.BookRequest, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "genres: should be for exaple: 123,456,789.")
			},
		},
		{
			name: "failure - invalid pages",
			url:  "authors=1,11&genres=1&min-pages=0&max-pages=10001&min-year=2001&max-year=2100&limit=5",
			assert: func(resp models.BookRequest, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "max-pages: should be between 1 and 10000; min-pages: should be between 1 and 10000.")
			},
		},
		{
			name: "failure - invalid year",
			url:  "authors=1,11&genres=1&min-pages=100&max-pages=300&min-year=1799&max-year=2101&limit=5",
			assert: func(resp models.BookRequest, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "max-year: should be between 1800 and 2100; min-year: should be between 1800 and 2100.")
			},
		},
		{
			name: "failure - invalid limit",
			url:  "authors=1,11&genres=1&min-pages=100&max-pages=300&min-year=2001&max-year=2100&limit=0",
			assert: func(resp models.BookRequest, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), "limit: should be between 1 and 1000.")
			},
		},
	}

	for _, c := range cases {
		req := http.Request{
			URL: &url.URL{RawQuery: c.url},
		}
		bookReq := translators.ToBooksRequest(&req)
		err := bookReq.Validate()
		c.assert(bookReq, err)
	}
}
