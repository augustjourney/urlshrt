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

var UrlsInMemory []storage.URL

func (r *Repo) Create(short string, original string) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Log.Error("Could not create uuid ", err)
		panic(err)
	}
	url := storage.URL{
		UUID:     uuid.String(),
		Short:    short,
		Original: original,
	}

	urls := r.GetAll()
	urls = append(urls, url)

	data, err := json.Marshal(&urls)
	if err != nil {
		logger.Log.Error("Could not marshal json urls ", err)
		panic(err)
	}
	file, err := os.OpenFile(r.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Could not open file to write urls ", err)
		panic(err)
	}
	defer file.Close()
	file.Write(data)
}

func (r *Repo) GetAll() []storage.URL {
	file, err := os.OpenFile(r.fileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log.Error("Could not open file to read all urls ", err)
		panic(err)
	}

	defer file.Close()

	var urls []storage.URL

	data, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Error("Could not read all urls ", err)
		panic(err)
	}

	err = json.Unmarshal(data, &urls)
	if err != nil {
		// if err is unexpected end of JSON input
		// it means json file is empty
		if !strings.Contains(err.Error(), "unexpected end of JSON input") {
			logger.Log.Error("Could not unmarshal all urls ", err)
			panic(err)
		}
	}

	return urls
}

func (r *Repo) Get(short string) *storage.URL {

	urls := r.GetAll()

	var url storage.URL

	for i := 0; i < len(urls); i++ {
		if urls[i].Short == short {
			url = urls[i]
			break
		}
	}

	return &url
}

func New(config *config.Config) Repo {
	return Repo{
		fileStoragePath: config.FileStoragePath,
	}
}
