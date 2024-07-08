// модуль отвечает за сохранение данных о ссылках в оперативной памяти.
package inmemory

import (
	"context"

	"github.com/augustjourney/urlshrt/internal/storage"
)

// репозиторий с методами хранилища
type Repo struct{}

// слайс для хранения ссылок в памяти
var UrlsInMemory []storage.URL

// сохраняет ссылку в хранилище
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

// сохраняет множество ссылок в хранилище
func (r *Repo) CreateBatch(ctx context.Context, urls []storage.URL) error {
	UrlsInMemory = append(UrlsInMemory, urls...)
	return nil
}

// получает экземпляр ссылки по короткой
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

// получает внутренню статистику: количество сохранненых ссылок и количество пользователей
func (r *Repo) GetStats(ctx context.Context) (storage.Stats, error) {
	var urlsCount int
	users := make(map[string]bool)
	for i := 0; i < len(UrlsInMemory); i++ {
		if UrlsInMemory[i].UserUUID != "" {
			users[UrlsInMemory[i].UserUUID] = true
		}
		if UrlsInMemory[i].Short != "" {
			urlsCount++
		}
	}
	usersCount := len(users)
	return storage.Stats{
		UrlsCount:  urlsCount,
		UsersCount: usersCount,
	}, nil
}

// получает ссылки пользователя
func (r *Repo) GetByUserUUID(ctx context.Context, userUUID string) (*[]storage.URL, error) {
	var urls []storage.URL

	for i := 0; i < len(UrlsInMemory); i++ {
		if UrlsInMemory[i].UserUUID == userUUID && !UrlsInMemory[i].IsDeleted {
			urls = append(urls, UrlsInMemory[i])
		}
	}

	return &urls, nil
}

// удаляет ссылку
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

// получает экземпляр ссылки по оригинальной
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

// создает новый экземпляр inmemory-репозитория
func New() *Repo {
	UrlsInMemory = make([]storage.URL, 100000)
	return &Repo{}
}
