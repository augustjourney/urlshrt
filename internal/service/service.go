package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/google/uuid"
)

var errNotFound = errors.New("url not found")
var errInternalError = errors.New("internal error")

type Service struct {
	repo   storage.IRepo
	config *config.Config
}

type IService interface {
	Shorten(originalURL string) (string, error)
	FindOriginal(short string) (string, error)
	ShortenBatch(batchURLs []BatchURL) ([]BatchResultURL, error)
}

type BatchURL struct {
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

type BatchResultURL struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

func (s *Service) generateID() (string, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		logger.Log.Error("Could not create uuid ", err)
		return "", err
	}

	return uuid.String(), nil
}

func (s *Service) hashURL(url string) string {
	hash := sha256.New()
	io.WriteString(hash, url)
	return fmt.Sprintf("%x", hash.Sum(nil))[:10]
}

func (s *Service) Shorten(originalURL string) (string, error) {
	short := s.hashURL(originalURL)
	uuid, err := s.generateID()
	if err != nil {
		return "", errInternalError
	}
	err = s.repo.Create(context.TODO(), storage.URL{
		UUID:     uuid,
		Short:    short,
		Original: originalURL,
	})

	if err != nil {
		return "", errInternalError
	}

	return s.config.BaseURL + "/" + short, nil
}

func (s *Service) ShortenBatch(batchURLs []BatchURL) ([]BatchResultURL, error) {
	var urls []storage.URL
	var result []BatchResultURL

	for _, url := range batchURLs {

		// Если url пришел без correlation_id
		// То не обрабатываем
		// Хотя наподумать — можно будет какой-то рандомный ID присваивать тогда
		if url.CorrelationID == "" {
			continue
		}

		short := s.hashURL(url.OriginalURL)

		uuid, err := s.generateID()

		if err != nil {
			return nil, err
		}

		urls = append(urls, storage.URL{
			Short:    short,
			Original: url.OriginalURL,
			UUID:     uuid,
		})

		result = append(result, BatchResultURL{
			CorrelationID: url.CorrelationID,
			ShortURL:      s.config.BaseURL + "/" + short,
		})
	}

	err := s.repo.CreateBatch(context.TODO(), urls)

	if err != nil {
		logger.Log.Error(err)
		return result, errInternalError
	}

	return result, nil
}

func (s *Service) FindOriginal(short string) (string, error) {
	url, err := s.repo.Get(context.TODO(), short)
	if err != nil {
		return "", errInternalError
	}
	if url == nil {
		return "", errNotFound
	}
	return url.Original, nil
}

func New(repo storage.IRepo, config *config.Config) Service {
	return Service{
		repo:   repo,
		config: config,
	}
}
