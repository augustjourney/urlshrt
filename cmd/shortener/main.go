package main

import (
	"net/http"

	"github.com/augustjourney/urlshrt/internal/controller"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage/inmemory"
)

func main() {
	mux := http.NewServeMux()

	repo := inmemory.New()
	service := service.New(&repo)
	c := controller.New(&service)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodPost {
			c.CreateURL(w, r)
			return
		} else if r.URL.Path != "/" && r.Method == http.MethodGet {
			c.GetURL(w, r)
			return
		} else {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	})

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}
