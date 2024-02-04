package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "http://localhost:8080"

type URL struct {
	Short    string
	Original string
}

// Storage

type Repo struct{}

var repo Repo

var UrlsInMemory []URL

func (r *Repo) Create(short string, original string) {
	url := URL{
		Short:    short,
		Original: original,
	}
	UrlsInMemory = append(UrlsInMemory, url)
}

func (r *Repo) Get(short string) *URL {

	for i := 0; i < len(UrlsInMemory); i++ {
		if UrlsInMemory[i].Short == short {
			return &UrlsInMemory[i]
		}
	}

	return nil
}

// Service, main logic

type Service struct {
	repo Repo
}

var service Service

func (s *Service) Shorten(originalURL string) string {
	short := "EwHXdJfB"
	s.repo.Create(short, originalURL)
	return BaseURL + "/" + short
}

func (s *Service) FindOriginal(short string) (string, error) {
	url := s.repo.Get(short)
	if url == nil {
		return "", errors.New("Url not found")
	}
	return url.Original, nil
}

// Http Controller

type Controller struct {
	service Service
}

func (c *Controller) urlHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == http.MethodPost {
		// Getting url from text plain body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		originalURL := string(body)

		// Make a short url
		short := c.service.Shorten(originalURL)

		// Response
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(short))
		return
	} else if r.URL.Path != "/" && r.Method == http.MethodGet {

		// Parse short url
		short := r.URL.Path[1:]

		// Find original
		originalURL, err := c.service.FindOriginal(short)

		fmt.Println("originalURL", originalURL)

		if err != nil {
			// Response
			fmt.Println(err)
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Response
		w.Header().Add("Location", originalURL)
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(originalURL))
		return
	} else {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	UrlsInMemory = make([]URL, 0)

	repo = Repo{}

	service = Service{
		repo: repo,
	}

	c := Controller{
		service: service,
	}

	mux.HandleFunc("/", c.urlHandler)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
