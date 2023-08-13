package dbwork

import (
	"database/sql"
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

// user=postgres password=490Sutud dbname=link-shortener sslmode=disable

func CreateTables(db *sql.DB, createTableQuery string) error {
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	log.Println("Tables created")
	return nil
}

func AddURL(db *sql.DB, shortURL, originalURL string) error {
	_, err := db.Exec("INSERT INTO urls (shortURL, originalURL) VALUES ($1, $2)", shortURL, originalURL)
	if err != nil {
		return err
	}

	log.Println("URL added to database")

	return nil
}

func GetOriginalURL(db *sql.DB, shortURL string) (string, error) {
	var originalURL string
	err := db.QueryRow("SELECT originalURL FROM urls WHERE shortURL = $1", shortURL).Scan(&originalURL)
	if err != nil {
		return "", err
	}

	log.Println("Original URL retrieved from database")

	return originalURL, nil
}
