package dbwork

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func CreateTestTable(db *sql.DB) error {
	const createTableQuery = `CREATE TEMPORARY TABLE test (id SERIAL PRIMARY KEY, shortURL TEXT UNIQUE, originalURL TEXT)`
	_, err := db.Exec(createTableQuery)
	return err
}
func DropTestTables(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS test")
	return err
}
func TestDBWork(t *testing.T) {
	// Подключение к тестовой базе данных
	connStr := "user=postgres password=490Sutud dbname=link-shortener sslmode=disable"
	testDB, err := sql.Open("postgres", connStr)
	assert.NoError(t, err, "Failed to connect to test database")
	defer testDB.Close()

	// Создание тестовых таблиц
	err = CreateTestTable(testDB)
	assert.NoError(t, err, "Failed to create table")

	// Добавление и получение URL
	// uuid := uuid.New().String()
	shortURL := uuid.New().String()
	originalURL := uuid.New().String()
	err = AddURL(testDB, shortURL, originalURL)
	assert.NoError(t, err, "Failed to add URL")

	retrievedURL, err := GetOriginalURL(testDB, shortURL)
	assert.NoError(t, err, "Failed to get original URL")
	assert.Equal(t, originalURL, retrievedURL, "Retrieved URL does not match")

	// Очистка тестовой базы данных (нужно предварительно удалить таблицы)
	// _, err = testDB.Exec("DROP TABLE IF EXISTS test")
	// assert.NoError(t, err, "Failed to drop tables")
}

func TestConnect(t *testing.T) {
	connStr := "user=postgres password=490Sutud dbname=testDB sslmode=disable" // Замените на параметры своей тестовой базы данных
	db, err := Connect(connStr)
	assert.NoError(t, err, "Failed to connect to test database")
	defer db.Close()

	// Проверка соединения с базой данных
	err = db.Ping()
	assert.NoError(t, err, "Failed to ping database")
}
