package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/storage"
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
}

func (s *Service) Shorten(originalURL string) (string, error) {
	hash := sha256.New()
	io.WriteString(hash, originalURL)
	short := fmt.Sprintf("%x", hash.Sum(nil))[:10]
	err := s.repo.Create(short, originalURL)
	if err != nil {
		return "", errInternalError
	}
	return s.config.BaseURL + "/" + short, nil
}

func (s *Service) FindOriginal(short string) (string, error) {
	url, err := s.repo.Get(short)
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
