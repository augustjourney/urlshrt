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

type Service struct {
	repo   storage.IRepo
	config *config.Config
}

type IService interface {
	Shorten(originalURL string) string
	FindOriginal(short string) (string, error)
}

func (s *Service) Shorten(originalURL string) string {
	hash := sha256.New()
	io.WriteString(hash, originalURL)
	short := fmt.Sprintf("%x", hash.Sum(nil))[:10]
	s.repo.Create(short, originalURL)
	return s.config.BaseURL + "/" + short
}

func (s *Service) FindOriginal(short string) (string, error) {
	url := s.repo.Get(short)
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
