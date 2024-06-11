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
	"strconv"
	"strings"
	"testing"

	"github.com/augustjourney/urlshrt/internal/storage/inmemory"
	"github.com/google/uuid"

	"github.com/augustjourney/urlshrt/internal/app"
	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAppInstance() (*fiber.App, storage.IRepo, service.Service) {
	config := config.New()
	logger.New()

	repo := inmemory.New()
	service := service.New(repo, config)
	controller := New(&service)

	app := app.New(&controller, nil)

	return app, repo, service
}

func TestGetURL(t *testing.T) {

	app, repo, _ := newAppInstance()

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

func BenchmarkGetURL(b *testing.B) {
	app, repo, _ := newAppInstance()

	url := storage.URL{
		UUID:     "some-uuid-1",
		Short:    "shrturl1",
		Original: "http://google.com",
	}

	repo.Create(context.TODO(), url)

	b.Run("create url", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			request := httptest.NewRequest(http.MethodGet, "/"+url.Short, nil)
			result, err := app.Test(request, 100)
			require.NoError(b, err)
			result.Body.Close()
		}
	})
}

func TestCreateURL(t *testing.T) {
	app, _, _ := newAppInstance()

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
			name: "URL created with conflict",
			want: want{
				code:        http.StatusConflict,
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
			result, err := app.Test(request, 100)
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

func BenchmarkCreateURL(b *testing.B) {
	app, _, _ := newAppInstance()

	b.Run("create url", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			originalURL, err := uuid.NewRandom()
			require.NoError(b, err)
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(originalURL.String())))
			result, err := app.Test(request, 100)
			require.NoError(b, err)
			result.Body.Close()
		}
	})
}

func TestApiCreateURL(t *testing.T) {
	app, _, _ := newAppInstance()

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
			name: "URL created with conflict",
			want: want{
				code:        http.StatusConflict,
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

			result, err := app.Test(request, 100)
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

func BenchmarkApiCreateURL(b *testing.B) {
	app, _, _ := newAppInstance()

	b.Run("create url", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			originalURL, err := uuid.NewRandom()
			require.NoError(b, err)
			body, _ := json.Marshal(APICreateURLBody{
				URL: originalURL.String(),
			})
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
			result, err := app.Test(request, 100)
			require.NoError(b, err)
			result.Body.Close()
		}
	})
}

func TestApiCreateURLBatch(t *testing.T) {
	app, _, _ := newAppInstance()

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
			body: `[{
					"original_url": "http://yandex.ru/123123",
					"correlation_id": "1"
				},
					{
					"original_url": "http://vk.com/12322222",
					"correlation_id": "2"
				},
					{
					"original_url": "http://ya.ru/123000000",
					"correlation_id": "3"
				}]`,
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
				},]`,
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

			result, err := app.Test(request, 100)
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

func TestApiDeleteBatch(t *testing.T) {
	app, _, urlsService := newAppInstance()

	// Нужно создать какое-то количество урлов
	// Которые будут удалены
	// По ним будем проверять статус удаления
	// То есть будем проходить по ним и проверять
	// Действительно ли был удален и статус 409

	// Также создать какое-то количество урлов
	// Которые не будут удалены
	// По ним должен быть статус 307

	// Для создания всех этих урлов можно использовать batch create

	var batch1 []service.BatchURL
	var batch2 []service.BatchURL

	for i := 0; i < 1000; i++ {
		uuid, _ := urlsService.GenerateID()
		batch1 = append(batch1, service.BatchURL{
			CorrelationID: strconv.Itoa(i),
			OriginalURL:   fmt.Sprintf("http://random-url-delete/%s%d", uuid, i),
		})
	}

	for i := 0; i < 500; i++ {
		uuid, _ := urlsService.GenerateID()
		batch2 = append(batch2, service.BatchURL{
			CorrelationID: strconv.Itoa(i),
			OriginalURL:   fmt.Sprintf("http://random-url-store/%s%d", uuid, i),
		})
	}

	userID := "user-123"

	urlsToStore, err := urlsService.ShortenBatch(batch2, userID)
	assert.NoError(t, err)

	urlsToDelete, err := urlsService.ShortenBatch(batch1, userID)
	assert.NoError(t, err)

	// Сейчас созданы все необходимые урлы
	// Можно тестировать, удалять и смотреть на результаты

	url := "/api/user/urls"

	var shortUrls []string

	for _, url := range urlsToDelete {
		// В ответе от сервиса short url будем вместе с хостом
		// поэтому отделяем по последнему слэшу
		lastSlash := strings.LastIndex(url.ShortURL, "/")
		shortUrls = append(shortUrls, url.ShortURL[lastSlash+1:])
	}

	requestBody, err := json.Marshal(shortUrls)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodDelete, url, bytes.NewReader(requestBody))

	request.Header.Set("authorization", userID)

	request.Header.Set("Content-Type", "application/json")

	result, err := app.Test(request, 60000)
	require.NoError(t, err)

	assert.Equal(t, http.StatusAccepted, result.StatusCode)

	for _, url := range urlsToStore {
		request := httptest.NewRequest(http.MethodGet, url.ShortURL, nil)

		require.NoError(t, err)
		res, err := app.Test(request, 500)
		require.NoError(t, err)
		require.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		res.Body.Close()
	}

	for _, url := range urlsToDelete {
		request := httptest.NewRequest(http.MethodGet, url.ShortURL, nil)

		require.NoError(t, err)
		res, err := app.Test(request, 100)
		require.NoError(t, err)
		require.Equal(t, http.StatusGone, res.StatusCode)
		res.Body.Close()
	}

	_ = result.Body.Close()

}
