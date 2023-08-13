package config

import (
	"os"
)

type Config struct {
	Address       string
	BaseURL       string
	FileStore     string
	DBConnect     string
	ServerAddrENV string `env:"SERVER_ADDRESS"`
	BaseURLENV    string `env:"BASE_URL"`
	FileENV       string `env:"FILE_STORAGE_PATH"`
	dbENV         string `env:"DATABASE_DSN"`
}

func GetBaseURL(cfg Config) string {
	return cfg.BaseURL
}

func GetConfig() Config {
	var conf Config
	if address := os.Getenv("SERVER_ADDRESS"); address != "" {
		conf.Address = address
	} else {
		conf.Address = "localhost:8080"
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		conf.BaseURL = baseURL
	} else {
		conf.BaseURL = "http://localhost:8080"
	}

	if fileName := os.Getenv("FILE_STORAGE_PATH"); fileName != "" {
		conf.FileStore = fileName
	} else {
		conf.FileStore = "short-url-db.json"
	}

	if dbConnect := os.Getenv("DATABASE_DSN"); dbConnect != "" {
		conf.DBConnect = dbConnect
	}
	// else {
	// 	conf.DBConnect = ""
	// }

	return conf
}
