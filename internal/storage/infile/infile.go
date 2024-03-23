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

type Repo struct {
	fileStoragePath string
}

/*
Логику добавления URL сделал такой:
- Создаем URL struct
- Получаем все текущие URLs с типом []storage.URL
- И в этот слайс добавляем новый урл
- И этот слайс перезаписываем в json

Мне кажется, это не совсем правильная логика.
Правильнее было бы добавлять один url в конец файла.
Но тогда json будет невалидный или я не понял, как это сделать.
Тогда самому нужно проверять запятые и конец файла.
Пока не разобрался с этим.
*/
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

func (r *Repo) GetByUserUUID(ctx context.Context, userUUID string) (*[]storage.URL, error) {
	var urls []storage.URL

	allURLs, err := r.GetAll(ctx)

	if err != nil {
		return &urls, nil
	}

	for i := 0; i < len(allURLs); i++ {
		if allURLs[i].UserUUID == userUUID {
			urls = append(urls, urls[i])
		}
	}

	return &urls, nil
}

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

func New(config *config.Config) *Repo {
	repo := Repo{
		fileStoragePath: config.FileStoragePath,
	}
	return &repo
}
