package config

import (
	"flag"
)

type Config struct {
 Address string
 BaseURL string
}

func GetConfig() *Config {
 address := flag.String("a", "localhost:8080", "address to run HTTP server")
 baseURL := flag.String("b", "http://localhost:8080/", "base URL for shortened URL")
 flag.Parse()

 return &Config{
  Address: *address,
  BaseURL: *baseURL,
 }
}