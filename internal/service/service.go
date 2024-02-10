package service

import (
	"errors"

	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/storage"
)

type Service struct {
	repo   storage.IRepo
	config *config.Config
}

type IService interface {
	Shorten(originalURL string) string
	FindOriginal(short string) (string, error)
}

func (s *Service) Shorten(originalURL string) string {
	short := "EwHXdJfB"
	s.repo.Create(short, originalURL)
	return s.config.BaseURL + "/" + short
}

func (s *Service) FindOriginal(short string) (string, error) {
	url := s.repo.Get(short)
	if url == nil {
		return "", errors.New("url not found")
	}
	return url.Original, nil
}

func New(repo storage.IRepo, config *config.Config) Service {
	return Service{
		repo:   repo,
		config: config,
	}
}
