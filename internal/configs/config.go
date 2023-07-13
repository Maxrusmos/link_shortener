package config

import (
	"os"
)

type Config struct {
	Address       string
	BaseURL       string
	ServerAddrENV string `env:"SERVER_ADDRESS"`
	BaseURLENV    string `env:"BASE_URL"`
}

func GetBaseURL(cfg Config) string {
	return cfg.BaseURL
}

func GetConfig() Config {
	var conf Config

	// conf.Address = "localhost:8080"
	// conf.BaseURL = "http://localhost:8080"

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

	// fmt.Println(conf)
	return conf
}
