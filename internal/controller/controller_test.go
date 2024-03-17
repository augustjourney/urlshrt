package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/augustjourney/urlshrt/internal/app"
	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/augustjourney/urlshrt/internal/storage/inmemory"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAppInstance() (*fiber.App, storage.IRepo) {
	config := config.New()
	logger.New()

	repo := inmemory.New()
	service := service.New(repo, config)
	controller := New(&service)

	app := app.New(&controller, nil)

	return app, repo
}

func TestGetURL(t *testing.T) {

	app, repo := newAppInstance()

	url1 := storage.URL{
		UUID:     "some-uuid-1",
		Short:    "shrturl1",
		Original: "http://google.com",
	}

	url2 := storage.URL{
		UUID:     "some-uuid-2",
		Short:    "shrturl2",
		Original: "http://yandex.ru",
	}

	repo.Create(context.TODO(), url1)
	repo.Create(context.TODO(), url2)

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
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
			shortURL: "3453",
		},
		{
			name:   "Found url",
			method: http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
				response:    url1.Original,
			},
			shortURL: url1.Short,
		},
		{
			name:   "Found url 2",
			method: http.MethodGet,
			want: want{
				code:        http.StatusTemporaryRedirect,
				contentType: "text/plain",
				response:    url2.Original,
			},
			shortURL: url2.Short,
		},
		{
			name:   "Method [PUT] not allowed",
			method: http.MethodPut,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
			shortURL: url2.Short,
		},
		{
			name:   "Method [DELETE] not allowed",
			method: http.MethodDelete,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
			shortURL: url2.Short,
		},
		{
			name:   "Method [POST] not allowed with shortURL",
			method: http.MethodDelete,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
			shortURL: url2.Short,
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
			res.Body.Close()
			if tt.method == http.MethodGet {
				assert.Equal(t, tt.want.response, res.Header.Get("Location"))
			}
		})
	}
}

func TestCreateURL(t *testing.T) {
	app, _ := newAppInstance()

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
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
			originalURL: "http://yandex.ru",
			method:      http.MethodPost,
		},
		{
			name: "Wront HTTP method",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
			originalURL: "http://yandex.ru",
			method:      http.MethodPut,
		},
		{
			name: "Empty body",
			want: want{
				code:        http.StatusBadRequest,
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
			result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			if tt.method == http.MethodPost && tt.originalURL != "" {
				assert.Equal(t, true, shortMatch)
			}
		})
	}
}

func TestApiCreateURL(t *testing.T) {
	app, _ := newAppInstance()

	type want struct {
		code        int
		contentType string
		result      bool
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
				code:        http.StatusCreated,
				contentType: "application/json",
				result:      true,
			},
			originalURL: "http://yandex.ru",
			method:      http.MethodPost,
		},
		{
			name: "Empty body",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "application/json",
				result:      false,
			},
			originalURL: "",
			method:      http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/shorten"

			body, _ := json.Marshal(APICreateURLBody{
				URL: tt.originalURL,
			})

			request := httptest.NewRequest(tt.method, url, bytes.NewReader(body))
			request.Header.Set("Content-Type", tt.want.contentType)

			result, err := app.Test(request, 1)
			require.NoError(t, err)

			var resultBody APICreateURLResult

			err = json.NewDecoder(result.Body).Decode(&resultBody)
			result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"), fmt.Sprintf("Content Type should be %s", tt.want.contentType))

			if tt.want.result {
				require.NoError(t, err)
				assert.Equal(t, true, resultBody.Result != "", "Result url should not be empty")
			}
		})
	}
}

func TestApiCreateURLBatch(t *testing.T) {
	app, _ := newAppInstance()

	type want struct {
		code             int
		contentType      string
		resultUrlsLength int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "Batch URL created",
			body: `[
				{
					"original_url": "http://yandex.ru/123",
					"correlation_id": "1"
				},
					{
					"original_url": "http://vk.com/123",
					"correlation_id": "2"
				},
					{
					"original_url": "http://ya.ru/123",
					"correlation_id": "3"
				}
			]`,
			want: want{
				contentType:      "application/json",
				code:             http.StatusCreated,
				resultUrlsLength: 3,
			},
		},
		{
			name: "Batch URL Bad Request — Wrong JSON body",
			body: `[
				{
					"original_url": "http://yandex.ru/123",
					"correlation_id": "1"
				},
					{
					"original_url": "http://vk.com/123",
					"correlation_id": "2"
				},
					{
					"
			]`,
			want: want{
				contentType:      "application/json",
				code:             http.StatusBadRequest,
				resultUrlsLength: 0,
			},
		},
		{
			name: "Batch URL Created — Only With Correlation ID",
			body: `[
				{
					"original_url": "http://yandex.ru/123",
					"correlation_id": "1"
				},
				{
					"original_url": "http://vk.com/123"
				}
			]`,
			want: want{
				contentType:      "application/json",
				code:             http.StatusCreated,
				resultUrlsLength: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/shorten/batch"

			request := httptest.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(tt.body)))

			request.Header.Set("Content-Type", tt.want.contentType)

			result, err := app.Test(request, 1)
			require.NoError(t, err)

			var resultBody []service.BatchResultURL

			err = json.NewDecoder(result.Body).Decode(&resultBody)
			result.Body.Close()

			assert.Equal(t, tt.want.code, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"), fmt.Sprintf("Content Type should be %s", tt.want.contentType))

			if result.StatusCode == http.StatusCreated {
				require.NoError(t, err)
				assert.Equal(t, tt.want.resultUrlsLength, len(resultBody))
			}
		})
	}
}
