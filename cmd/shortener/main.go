package main

import (
	"net/http"
)

func urlHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == http.MethodPost {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/EwHXdJfB"))
	} else if r.URL.Path != "/" && r.Method == http.MethodGet {
		w.Header().Add("Location", "https://practicum.yandex.ru/")
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("https://practicum.yandex.ru/"))
	} else {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
