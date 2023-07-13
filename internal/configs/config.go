package config

type Config struct {
	Address       string
	BaseURL       string
	ServerAddrENV string `env:"SERVER_ADDRESS"`
	BaseURLENV    string `env:"BASE_URL"`
}

func GetConfig() *Config {
	address := "localhost:8080"
	baseURL := "http://localhost:8080"
	serveraddrEnv := "SERVER_ADDRESS"
	baseurlEnv := "BASE_URL"

	return &Config{
		Address:       address,
		BaseURL:       baseURL,
		ServerAddrENV: serveraddrEnv,
		BaseURLENV:    baseurlEnv,
	}
}
