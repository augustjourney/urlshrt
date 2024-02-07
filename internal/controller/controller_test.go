package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURL(t *testing.T) {
	repo := inmemory.New()
	service := service.New(&repo)
	controller := New(&service)

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
			name:   "Found url",
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
		{
			name:   "Found URL — response",
			method: http.MethodGet,
			want: want{
				code:        307,
				contentType: "text/plain",
				response:    "http://yandex.ru",
			},
			shortURL: "123",
		},
		{
			name:   "Found URL — response 2",
			method: http.MethodGet,
			want: want{
				code:        307,
				contentType: "text/plain",
				response:    "http://google.com",
			},
			shortURL: "321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/" + tt.shortURL
			req, err := http.NewRequest(tt.method, url, nil)
			require.NoError(t, err)
			w := httptest.NewRecorder()
			controller.GetURL(w, req)
			res := w.Result()
			err = res.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, res.Header.Get("Location"))
		})
	}
}
