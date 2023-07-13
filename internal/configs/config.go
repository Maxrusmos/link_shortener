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

func GetConfig() Config {
	var conf Config

	// address := "localhost:8080"
	// baseURL := "http://localhost:8090"
	// serveraddrEnv := "SERVER_ADDRESS"
	// baseurlEnv := "BASE_URL"

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

	// // Проверка аргументов командной строки
	// flag.StringVar(&conf.Address, "a", conf.Address, "HTTP server address")
	// flag.StringVar(&conf.BaseURL, "b", conf.BaseURL, "Base address for shortened URL")
	// // flag.Parse()
	return conf

	// return &Config{
	// 	Address:       address,
	// 	BaseURL:       baseURL,
	// 	ServerAddrENV: serveraddrEnv,
	// 	BaseURLENV:    baseurlEnv,
	// }
}
