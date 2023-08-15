package storage

import (
	"testing"
)

func TestMapURLStorage(t *testing.T) {
	testStorage(t, NewMapURLStorage())
}

func TestFileURLStorage(t *testing.T) {
	testStorage(t, NewFileURLStorage("test_storage.json"))
}

// func TestDatabaseURLStorage(t *testing.T) {
// 	// Создаем временную базу данных для тестов
// 	db, err := sql.Open("postgres", "user=postgres password=490Sutud dbname=link-shorten sslmode=disable")
// 	if err != nil {
// 		t.Fatalf("Failed to create Postgres database: %v", err)
// 	}
// 	defer db.Close()

// 	// Создаем таблицу в тестовой базе данных
// 	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS testBD (
// 		id SERIAL PRIMARY KEY,
// 		shortURL TEXT UNIQUE,
// 		originalURL TEXT
// 	  )`)
// 	if err != nil {
// 		t.Fatalf("Failed to create test table: %v", err)
// 	}

// 	testStorage(t, NewDatabaseURLStorage(db))
// }

func testStorage(t *testing.T, storage URLStorage) {
	// Тест AddURL
	err := storage.AddURL("abc123", "http://example.com")
	if err != nil {
		t.Errorf("AddURL returned an error: %v", err)
	}

	// Тест GetURL
	url, err := storage.GetURL("abc123")
	if err != nil {
		t.Errorf("GetURL returned an error: %v", err)
	}
	if url != "http://example.com" {
		t.Errorf("Expected URL to be %s, but got %s", "http://example.com", url)
	}

	// Тест Ping
	// err = storage.Ping()
	// if err != nil {
	// 	t.Errorf("Ping returned an error: %v", err)
	// }
}

// Примечание: Вам может потребоваться адаптировать тесты в зависимости от ваших реализаций и логики работы функций.
