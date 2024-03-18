package storage

import (
	"context"
	"errors"
)

type URL struct {
	UUID     string `json:"uuid"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type IRepo interface {
	Create(ctx context.Context, url URL) error
	Get(ctx context.Context, short string) (*URL, error)
	GetByOriginal(ctx context.Context, original string) (*URL, error)
	CreateBatch(ctx context.Context, urls []URL) error
}

var ErrAlreadyExists = errors.New("URL already exists")
