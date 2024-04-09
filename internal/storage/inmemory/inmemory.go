package inmemory

import (
	"context"
	"github.com/augustjourney/urlshrt/internal/storage"
)

type Repo struct{}

var UrlsInMemory []storage.URL

func (r *Repo) Create(ctx context.Context, url storage.URL) error {
	foundURL, err := r.GetByOriginal(ctx, url.Original)
	if err != nil {
		return err
	}
	if foundURL.Short != "" {
		return storage.ErrAlreadyExists
	}
	UrlsInMemory = append(UrlsInMemory, url)
	return nil
}

func (r *Repo) CreateBatch(ctx context.Context, urls []storage.URL) error {
	UrlsInMemory = append(UrlsInMemory, urls...)
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

func (r *Repo) GetByUserUUID(ctx context.Context, userUUID string) (*[]storage.URL, error) {
	var urls []storage.URL

	for i := 0; i < len(UrlsInMemory); i++ {
		if UrlsInMemory[i].UserUUID == userUUID && !UrlsInMemory[i].IsDeleted {
			urls = append(urls, UrlsInMemory[i])
		}
	}

	return &urls, nil
}

func (r *Repo) Delete(ctx context.Context, shortURLs []string, userUUID string) error {

	shortUrlsMap := make(map[string]bool)

	for _, short := range shortURLs {
		shortUrlsMap[short] = true
	}

	for i := 0; i < len(UrlsInMemory); i++ {
		url := UrlsInMemory[i]
		_, ok := shortUrlsMap[url.Short]

		if userUUID == "" {
			UrlsInMemory[i].IsDeleted = true
		} else if url.UserUUID == userUUID && ok {
			UrlsInMemory[i].IsDeleted = true
		}
	}

	return nil
}

func (r *Repo) GetByOriginal(ctx context.Context, original string) (*storage.URL, error) {
	var url storage.URL
	for i := 0; i < len(UrlsInMemory); i++ {
		if UrlsInMemory[i].Original == original {
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
