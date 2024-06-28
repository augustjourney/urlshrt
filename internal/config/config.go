package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/augustjourney/urlshrt/internal/logger"
)

// Конфиг — хранит в себе настройки приложения
type Config struct {
	BaseURL         string `env:"BASE_URL" json:"base_url"`
	ServerAddress   string `env:"SERVER_ADDRESS" json:"server_address"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	EnableHTTPS     bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	CertPemPath     string `json:"-"`
	CertKeyPath     string `json:"-"`
	Config          string `env:"CONFIG" json:"-"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

var config Config

func parseJSONConfig(pathToConfigFile string, config *Config) {
	configFile, err := os.ReadFile(pathToConfigFile)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	err = json.Unmarshal(configFile, config)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

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
		flagServerAddress   = flag.String("a", "", "Server address on which server is running")
		flagBaseURL         = flag.String("b", "", "Base URL which short urls will be accessible")
		flagFileStoragePath = flag.String("f", "", "Path to file where urls data will be stored")
		flagDatabaseDSN     = flag.String("d", "", "Database DSN")
		flagEnableHTTPS     = flag.Bool("s", false, "Enable HTTPS")
		flagConfig          = flag.String("c", "", "Config in JSON")
		flagTrustedSubnet   = flag.String("t", "", "Trusted subnet")
	)

	flag.Parse()

	// Инициализация конфига с дефолтными значениями
	config = Config{
		CertPemPath:     defaults["certPemPath"],
		CertKeyPath:     defaults["certKeyPath"],
		Config:          *flagConfig,
		ServerAddress:   defaults["serverAddress"],
		BaseURL:         defaults["baseURL"],
		FileStoragePath: defaults["fileStoragePath"],
	}

	// Если указан путь до конфиг-файла из json, парсим его
	if *flagConfig != "" {
		parseJSONConfig(*flagConfig, &config)
	}

	// Берем переменные из флагов, если они есть
	if *flagServerAddress != "" {
		config.ServerAddress = *flagServerAddress
	}

	if *flagBaseURL != "" {
		config.BaseURL = *flagBaseURL
	}

	if *flagFileStoragePath != "" {
		config.FileStoragePath = *flagFileStoragePath
	}

	if *flagDatabaseDSN != "" {
		config.DatabaseDSN = *flagDatabaseDSN
	}

	if *flagEnableHTTPS {
		config.EnableHTTPS = *flagEnableHTTPS
	}

	if *flagTrustedSubnet != "" {
		config.TrustedSubnet = *flagTrustedSubnet
	}

	// Берем переменные из окружения
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

	if trustedSubnet := os.Getenv("TRUSTED_SUBNET"); trustedSubnet != "" {
		config.TrustedSubnet = trustedSubnet
	}

	if enableHTTPS := os.Getenv("ENABLE_HTTPS"); enableHTTPS != "" {
		enableHTTPS, err := strconv.ParseBool(os.Getenv("ENABLE_HTTPS"))
		if err == nil && enableHTTPS {
			config.EnableHTTPS = enableHTTPS
		}
	}

	return &config
}
