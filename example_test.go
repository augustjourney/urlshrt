package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/augustjourney/urlshrt/internal/controller"
	"github.com/augustjourney/urlshrt/internal/service"
)

const serverAddr = "http://localhost:8000"

func ExampleGetURL() {
	resp, err := http.Get(serverAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		resp.Body.Close()
	}()

	fmt.Println(body)
}

func ExampleCreateURL() {
	originalURL := "http://ya.ru"

	payload := bytes.NewReader([]byte(originalURL))
	contentType := "text/plain"

	resp, err := http.Post(serverAddr, contentType, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		resp.Body.Close()
	}()

	fmt.Println(body)
}

func ExampleAPICreateURL() {
	originalURL := "http://ya.ru"

	payload, err := json.Marshal(controller.APICreateURLBody{
		URL: originalURL,
	})

	contentType := "application/json"
	url := serverAddr + "/api/shorten"

	resp, err := http.Post(url, contentType, bytes.NewReader(payload))
	if err != nil {
		fmt.Println(err)
		return
	}


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		resp.Body.Close()
	}()

	fmt.Println(body)
}

func ExampleAPICreateURLBatch() {
	payload, err := json.Marshal([]service.BatchURL{
		{
			OriginalURL: "http://ya.ru",
			CorrelationID: "1",
		},
		{
			OriginalURL: "http://vk.com",
			CorrelationID: "2",
		},
	})

	contentType := "application/json"
	url := serverAddr + "/api/shorten/batch"

	resp, err := http.Post(url, contentType, bytes.NewReader(payload))

	if err != nil {
		fmt.Println(err)
		return
	}


	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		resp.Body.Close()
	}()

	fmt.Println(body)
}
