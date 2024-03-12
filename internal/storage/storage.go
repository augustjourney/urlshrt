package storage

import "context"

type URL struct {
	UUID     string `json:"uuid"`
	Short    string `json:"short_url"`
	Original string `json:"original_url"`
}

type IRepo interface {
	Create(ctx context.Context, short string, original string) error
	Get(ctx context.Context, short string) (*URL, error)
}
