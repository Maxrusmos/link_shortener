package dbwork

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Connected to database")
	return db, nil
}
func CreateTables(db *sql.DB, createTableQuery string) error {
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	log.Println("Tables created")
	return nil
}

func AddURL(db *sql.DB, shortURL, originalURL string, userID string) error {
	fmt.Println("add", shortURL, ":::", originalURL)
	_, err := db.Exec("INSERT INTO shortened_urls (short_url, original_url, user_id) VALUES ($1, $2, $3)", shortURL, originalURL, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	log.Println("URL added to database")

	return nil
}

func GetOriginalURL(db *sql.DB, shortURL string) (string, error) {
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM shortened_urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		return "error", err
	}
	return originalURL, nil
}
