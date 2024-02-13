package config

import (
	"flag"
	"os"
)

type Config struct {
	BaseURL       string `env:"BASE_URL"`
	ServerAddress string `env:"SERVER_ADDRESS"`
}

var config Config

func New() *Config {
	if config.BaseURL != "" {
		return &config
	}

	defaults := map[string]string{
		"baseURL":       "http://localhost:8080",
		"serverAddress": "localhost:8080",
	}

	var (
		flagServerAddress = flag.String("a", defaults["serverAddress"], "Server address on which server is running")
		flagBaseURL       = flag.String("b", defaults["baseURL"], "Base URL which short urls will be accessible")
	)

	flag.Parse()

	config = Config{
		ServerAddress: *flagServerAddress,
		BaseURL:       *flagBaseURL,
	}

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		config.ServerAddress = serverAddress
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	return &config
}
