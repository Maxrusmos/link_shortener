package dbwork

import (
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

func TestAddURL(t *testing.T) {
	testdb, err := Connect("user=postgres password=490Sutud dbname=testDB sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer testdb.Close()

	fmt.Println("ff")

	_, err = testdb.Exec("CREATE TABLE IF NOT EXISTS urls (id SERIAL PRIMARY KEY, shortURL TEXT UNIQUE, originalURL TEXT)")

	log.Println("Table created")
	if err != nil {
		t.Fatal(err)
	}

	shortURL := "a9b9f043"
	originalURL := "https://example.com"

	err = AddURL(testdb, shortURL, originalURL)
	if err != nil {
		t.Fatal(err)
	}

	retrievedURL, err := GetOriginalURL(testdb, shortURL)
	if err != nil {
		t.Fatal(err)
	}

	if retrievedURL != originalURL {
		t.Errorf("Expected retrieved URL to be %s, but got %s", originalURL, retrievedURL)
	}
}
