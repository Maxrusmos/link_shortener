package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBaseURL(t *testing.T) {
	conf := Config{
		BaseURL: "http://example.com",
	}
	baseURL := GetBaseURL(conf)
	assert.Equal(t, "http://example.com", baseURL, "Base URLs do not match")
}

func TestGetConfig(t *testing.T) {
	os.Setenv("SERVER_ADDRESS", "test_server")
	os.Setenv("BASE_URL", "test_base_url")
	os.Setenv("FILE_STORAGE_PATH", "test_file_path")
	os.Setenv("DATABASE_DSN", "test_db_dsn")

	conf := GetConfig()
	assert.Equal(t, "test_server", conf.Address, "Server addresses do not match")
	assert.Equal(t, "test_base_url", conf.BaseURL, "Base URLs do not match")
	assert.Equal(t, "test_file_path", conf.FileStore, "File storage paths do not match")
	assert.Equal(t, "test_db_dsn", conf.DBConnect, "DB connections do not match")

	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("DATABASE_DSN")
}
