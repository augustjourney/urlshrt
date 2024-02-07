package controller

import (
	"io"
	"net/http"

	"github.com/augustjourney/urlshrt/internal/service"
)

type Controller struct {
	service service.IService
}

func (c *Controller) CreateURL(w http.ResponseWriter, r *http.Request) {
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
}

func (c *Controller) GetURL(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse short url
	short := r.URL.Path[1:]

	// Find original
	originalURL, err := c.service.FindOriginal(short)

	if err != nil || originalURL == "" {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Response
	w.Header().Add("Location", originalURL)
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Write([]byte(originalURL))
}

func New(service service.IService) Controller {
	return Controller{
		service: service,
	}
}
