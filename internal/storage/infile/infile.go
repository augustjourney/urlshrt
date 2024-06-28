// модуль отвечает за сохранение данных о ссылках в файле.
package infile

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/storage"
)

// репозиторий с методами хранилища
type Repo struct {
	fileStoragePath string
}

// сохраняет ссылку в файл
func (r *Repo) Create(ctx context.Context, url storage.URL) error {

	foundURL, err := r.GetByOriginal(ctx, url.Original)
	if err != nil {
		return err
	}
	if foundURL.Short != "" {
		return storage.ErrAlreadyExists
	}

	urls, err := r.GetAll(ctx)
	if err != nil {
		return err
	}

	urls = append(urls, url)

	data, err := json.Marshal(&urls)
	if err != nil {
		logger.Log.Error("Could not marshal json urls ", err)
		return err
	}

	file, err := os.OpenFile(r.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Could not open file to write urls ", err)
		return err
	}

	defer file.Close()

	file.Write(data)

	return nil
}

// сохраняет множество ссылок в файл
func (r *Repo) CreateBatch(ctx context.Context, urls []storage.URL) error {

	currentURLs, err := r.GetAll(ctx)
	if err != nil {
		return err
	}

	currentURLs = append(currentURLs, urls...)

	data, err := json.Marshal(currentURLs)
	if err != nil {
		logger.Log.Error("Could not marshal json urls ", err)
		return err
	}

	file, err := os.OpenFile(r.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Could not open file to write urls ", err)
		return err
	}

	defer file.Close()

	file.Write(data)

	return nil
}

// получает все сохраненные ссылки
func (r *Repo) GetAll(ctx context.Context) ([]storage.URL, error) {
	file, err := os.OpenFile(r.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	var urls []storage.URL
	if err != nil {
		logger.Log.Error("Could not open file to read all urls ", err)
		return urls, err
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Error("Could not read all urls ", err)
		return urls, err
	}

	err = json.Unmarshal(data, &urls)
	if err != nil {
		// if err is unexpected end of JSON input
		// it means json file is empty
		if !strings.Contains(err.Error(), "unexpected end of JSON input") {
			logger.Log.Error("Could not unmarshal all urls ", err)
			return urls, err
		}
	}

	return urls, nil
}

// получает ссылки конкретного пользователя
func (r *Repo) GetByUserUUID(ctx context.Context, userUUID string) (*[]storage.URL, error) {
	var urls []storage.URL

	allURLs, err := r.GetAll(ctx)

	if err != nil {
		return &urls, nil
	}

	for i := 0; i < len(allURLs); i++ {
		if allURLs[i].UserUUID == userUUID && !allURLs[i].IsDeleted {
			urls = append(urls, allURLs[i])
		}
	}

	logger.Log.Info(allURLs)

	return &urls, nil
}

// удаляет ссылку
func (r *Repo) Delete(ctx context.Context, shortURLs []string, userUUID string) error {

	allURLs, err := r.GetAll(ctx)

	if err != nil {
		return nil
	}

	shortUrlsMap := make(map[string]bool)

	for _, short := range shortURLs {
		shortUrlsMap[short] = true
	}

	for i := 0; i < len(allURLs); i++ {
		url := allURLs[i]
		_, ok := shortUrlsMap[url.Short]
		if userUUID == "" {
			allURLs[i].IsDeleted = true
		} else if url.UserUUID == userUUID && ok {
			allURLs[i].IsDeleted = true
		}
	}

	data, err := json.Marshal(&allURLs)
	if err != nil {
		logger.Log.Error("Could not marshal json urls ", err)
		return err
	}

	file, err := os.OpenFile(r.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Could not open file to write urls ", err)
		return err
	}

	defer file.Close()

	file.Write(data)

	return nil
}

// получает внутренню статистику: количество сохранненых ссылок и количество пользователей
func (r *Repo) GetStats(ctx context.Context) (storage.Stats, error) {
	return storage.Stats{}, nil
}

// получает экземпляр ссылки по короткой
func (r *Repo) Get(ctx context.Context, short string) (*storage.URL, error) {

	var url storage.URL

	urls, err := r.GetAll(ctx)

	if err != nil {
		return &url, nil
	}

	for i := 0; i < len(urls); i++ {
		if urls[i].Short == short {
			url = urls[i]
			break
		}
	}

	return &url, nil
}

// получает экземпляр ссылки по оригинальной
func (r *Repo) GetByOriginal(ctx context.Context, original string) (*storage.URL, error) {

	var url storage.URL

	urls, err := r.GetAll(ctx)

	if err != nil {
		return &url, nil
	}

	for i := 0; i < len(urls); i++ {
		if urls[i].Original == original {
			url = urls[i]
			break
		}
	}

	return &url, nil
}

// создает новый экземпляр infile-репозитория
func New(config *config.Config) *Repo {
	repo := Repo{
		fileStoragePath: config.FileStoragePath,
	}
	return &repo
}
