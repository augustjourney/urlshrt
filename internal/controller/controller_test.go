package controller

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/augustjourney/urlshrt/internal/app"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURL(t *testing.T) {
	repo := inmemory.New()
	service := service.New(&repo)
	controller := New(&service)

	app := app.New(&controller)

	repo.Create("321", "http://google.com")
	repo.Create("123", "http://yandex.ru")

	type want struct {
		code        int
		contentType string
		response    string
	}

	tests := []struct {
		name     string
		want     want
		method   string
		shortURL string
	}{
		{
			name:   "Not found url",
			method: http.MethodGet,
			want: want{
				code:        400,
				contentType: "text/plain",
			},
			shortURL: "3453",
		},
		{
			name:   "Found url",
			method: http.MethodGet,
			want: want{
				code:        307,
				contentType: "text/plain",
				response:    "http://yandex.ru",
			},
			shortURL: "123",
		},
		{
			name:   "Found url 2",
			method: http.MethodGet,
			want: want{
				code:        307,
				contentType: "text/plain",
				response:    "http://google.com",
			},
			shortURL: "321",
		},
		{
			name:   "Method [PUT] not allowed",
			method: http.MethodPut,
			want: want{
				code:        400,
				contentType: "text/plain",
			},
			shortURL: "321",
		},
		{
			name:   "Method [DELETE] not allowed",
			method: http.MethodDelete,
			want: want{
				code:        400,
				contentType: "text/plain",
			},
			shortURL: "321",
		},
		{
			name:   "Method [POST] not allowed with shortURL",
			method: http.MethodDelete,
			want: want{
				code:        400,
				contentType: "text/plain",
			},
			shortURL: "321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/" + tt.shortURL
			req, err := http.NewRequest(tt.method, url, nil)
			require.NoError(t, err)
			res, err := app.Test(req, 1)
			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)
			if tt.method == http.MethodGet {
				assert.Equal(t, tt.want.response, res.Header.Get("Location"))
			}
		})
	}
}

func TestCreateURL(t *testing.T) {
	repo := inmemory.New()
	service := service.New(&repo)
	controller := New(&service)

	app := app.New(&controller)

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name        string
		want        want
		originalURL string
		method      string
	}{
		{
			name: "URL created",
			want: want{
				code:        201,
				contentType: "text/plain",
			},
			originalURL: "http://yandex.ru",
			method:      http.MethodPost,
		},
		{
			name: "Wront HTTP method",
			want: want{
				code:        400,
				contentType: "text/plain",
			},
			originalURL: "http://yandex.ru",
			method:      http.MethodPut,
		},
		{
			name: "Empty body",
			want: want{
				code:        400,
				contentType: "text/plain",
			},
			originalURL: "",
			method:      http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/"
			request := httptest.NewRequest(tt.method, url, bytes.NewReader([]byte(tt.originalURL)))
			result, err := app.Test(request, 1)
			require.NoError(t, err)
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			shortMatch, err := regexp.Match(`\/\w+$`, resultBody)
			require.NoError(t, err)
			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			if tt.method == http.MethodPost && tt.originalURL != "" {
				assert.Equal(t, true, shortMatch)
			}
			result.Body.Close()
		})
	}
}
