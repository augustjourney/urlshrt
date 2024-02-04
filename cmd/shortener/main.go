package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const BASE_URL = "http://localhost:8080"

type Url struct {
	Short    string
	Original string
}

// Storage

type Repo struct{}

var repo Repo

var URLS_IN_MEMORY []Url

func (r *Repo) Create(short string, original string) {
	url := Url{
		Short:    short,
		Original: original,
	}
	URLS_IN_MEMORY = append(URLS_IN_MEMORY, url)
}

func (r *Repo) Get(short string) *Url {

	for i := 0; i < len(URLS_IN_MEMORY); i++ {
		if URLS_IN_MEMORY[i].Short == short {
			return &URLS_IN_MEMORY[i]
		}
	}

	return nil
}

// Service, main logic

type Service struct {
	repo Repo
}

var service Service

func (s *Service) Shorten(originalUrl string) string {
	short := "EwHXdJfB"
	s.repo.Create(short, originalUrl)
	return BASE_URL + "/" + short
}

func (s *Service) FindOriginal(short string) (error, string) {
	url := s.repo.Get(short)
	if url == nil {
		return errors.New("Url not found"), ""
	}
	return nil, url.Original
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
		originalUrl := string(body)

		// Make a short url
		short := c.service.Shorten(originalUrl)

		// Response
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(short))
		return
	} else if r.URL.Path != "/" && r.Method == http.MethodGet {

		// Parse short url
		short := r.URL.Path[1:]

		// Find original
		err, originalUrl := c.service.FindOriginal(short)

		fmt.Println("originalUrl", originalUrl)

		if err != nil {
			// Response
			fmt.Println(err)
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Response
		w.Header().Add("Location", originalUrl)
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte(originalUrl))
		return
	} else {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	URLS_IN_MEMORY = make([]Url, 0)

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
