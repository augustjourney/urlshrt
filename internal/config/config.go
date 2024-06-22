package config

import (
	"flag"
	"os"
	"strconv"
)

// Конфиг — хранит в себе настройки приложения
type Config struct {
	BaseURL         string `env:"BASE_URL"`
	ServerAddress   string `env:"SERVER_ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS"`
	CertPemPath     string
	CertKeyPath     string
}

var config Config

// Создает экземпляр конфига
func New() *Config {
	if config.BaseURL != "" {
		return &config
	}

	defaults := map[string]string{
		"baseURL":         "http://localhost:8080",
		"serverAddress":   "localhost:8080",
		"fileStoragePath": "/tmp/short-url-db.json",
		"certPemPath":     "certs/cert.pem",
		"certKeyPath":     "certs/cert.key",
	}

	var (
		flagServerAddress = flag.String("a", defaults["serverAddress"], "Server address on which server is running")
		flagBaseURL       = flag.String("b", defaults["baseURL"], "Base URL which short urls will be accessible")
		fileStoragePath   = flag.String("f", defaults["fileStoragePath"], "Path to file where urls data will be stored")
		flagDatabaseDSN   = flag.String("d", "", "Database DSN")
		flagEnableHTTPS   = flag.Bool("s", false, "Enable HTTPS")
	)

	flag.Parse()

	config = Config{
		ServerAddress:   *flagServerAddress,
		BaseURL:         *flagBaseURL,
		FileStoragePath: *fileStoragePath,
		DatabaseDSN:     *flagDatabaseDSN,
		EnableHTTPS:     *flagEnableHTTPS,
		CertPemPath:     defaults["certPemPath"],
		CertKeyPath:     defaults["certKeyPath"],
	}

	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		config.ServerAddress = serverAddress
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	if fileStoragePath := os.Getenv("FILE_STORAGE_PATH"); fileStoragePath != "" {
		config.FileStoragePath = fileStoragePath
	}

	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		config.DatabaseDSN = databaseDSN
	}

	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		config.DatabaseDSN = databaseDSN
	}

	if enableHTTPS := os.Getenv("ENABLE_HTTPS"); enableHTTPS != "" {
		enableHTTPS, err := strconv.ParseBool(os.Getenv("ENABLE_HTTPS"))
		if err == nil && enableHTTPS {
			config.EnableHTTPS = enableHTTPS
		}
	}

	return &config
}
