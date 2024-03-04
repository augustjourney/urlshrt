package infile

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/google/uuid"
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
func (r *Repo) Create(short string, original string) error {
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

	urls, err := r.GetAll()
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

func (r *Repo) GetAll() ([]storage.URL, error) {
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

func (r *Repo) Get(short string) (*storage.URL, error) {

	var url storage.URL

	urls, err := r.GetAll()

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

func New(config *config.Config) Repo {
	return Repo{
		fileStoragePath: config.FileStoragePath,
	}
}
