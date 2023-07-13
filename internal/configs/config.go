package config

type Config struct {
	Address string
	BaseURL string
}

func GetConfig() *Config {
	address := "localhost:8080"
	baseURL := "http://localhost:8080"

	return &Config{
		Address: address,
		BaseURL: baseURL,
	}
}
