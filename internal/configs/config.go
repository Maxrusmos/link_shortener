package config

import (
	"link_shortener/internal/dbwork"
	"link_shortener/internal/storage"
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

func GetStorageURL(conf Config) storage.URLStorage {
	if conf.DBConnect == "" && conf.FileStore == "" && os.Getenv("FILE_STORAGE_PATH") == "" {
		return storage.NewMapURLStorage()
	}

	if conf.FileStore != "" || os.Getenv("FILE_STORAGE_PATH") != "" {
		return storage.NewFileURLStorage(conf.FileStore)
	}

	db, err := dbwork.Connect(conf.DBConnect)
	if err != nil {
		panic(err)
	}
	err = dbwork.CreateTables(db, `CREATE TABLE IF NOT EXISTS shortened_urls  (
        id SERIAL PRIMARY KEY,
        short_url VARCHAR(50) NOT NULL,
        original_url TEXT NOT NULL,
        UNIQUE (original_url)
      )`)
	if err != nil {
		panic(err)
	}
	return storage.NewDatabaseURLStorage(db)
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
	} else {
		conf.DBConnect = "user=postgres password=490Sutud dbname=link-shortener sslmode=disable"
	}

	return conf
}
