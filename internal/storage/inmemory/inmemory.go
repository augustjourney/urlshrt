package inmemory

import (
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/google/uuid"
)

type Repo struct{}

var UrlsInMemory []storage.URL

func (r *Repo) Create(short string, original string) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	url := storage.URL{
		UUID:     uuid.String(),
		Short:    short,
		Original: original,
	}
	UrlsInMemory = append(UrlsInMemory, url)
}

func (r *Repo) Get(short string) *storage.URL {
	var url storage.URL
	for i := 0; i < len(UrlsInMemory); i++ {
		if UrlsInMemory[i].Short == short {
			url = UrlsInMemory[i]
			break
		}
	}

	return &url
}

func New() Repo {
	UrlsInMemory = make([]storage.URL, 0)
	return Repo{}
}
