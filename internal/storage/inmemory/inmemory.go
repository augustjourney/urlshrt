package inmemory

import (
	"context"

	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/google/uuid"
)

type Repo struct{}

var UrlsInMemory []storage.URL

func (r *Repo) Create(ctx context.Context, short string, original string) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Log.Error("Could not create uuid ", err)
		return err
	}
	url := storage.URL{
		UUID:     uuid.String(),
		Short:    short,
		Original: original,
	}
	UrlsInMemory = append(UrlsInMemory, url)
	return nil
}

func (r *Repo) Get(ctx context.Context, short string) (*storage.URL, error) {
	var url storage.URL
	for i := 0; i < len(UrlsInMemory); i++ {
		if UrlsInMemory[i].Short == short {
			url = UrlsInMemory[i]
			break
		}
	}

	return &url, nil
}

func New() *Repo {
	UrlsInMemory = make([]storage.URL, 0)
	return &Repo{}
}
