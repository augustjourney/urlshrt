package config

import "flag"

type Config struct {
	BaseURL string
	Port    string
}

var config Config

func New() *Config {
	if config.BaseURL != "" {
		return &config
	}

	defaults := map[string]string{
		"baseURL": "http://localhost:8080",
		"port":    "8080",
	}

	var (
		port    = flag.String("a", defaults["port"], "Port on which server is running")
		baseURL = flag.String("b", defaults["baseURL"], "Base URL which short urls will be accessible")
	)

	flag.Parse()

	config = Config{
		Port:    *port,
		BaseURL: *baseURL,
	}

	return &config
}
