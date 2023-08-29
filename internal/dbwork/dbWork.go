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

func DeleteFromDB(db *sql.DB, urlsToDelete []string, userID string) error {
	stmt, err := db.Prepare("UPDATE shortened_urls SET deleted_flag = true WHERE short_url = $1 AND user_id = $2")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	for _, shortURL := range urlsToDelete {
		_, err := stmt.Exec(shortURL, userID)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	// deleteStmt, err := db.Prepare("DELETE FROM shortened_urls WHERE deleted_flag = true AND user_id = $1")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	// defer deleteStmt.Close()

	// _, err = deleteStmt.Exec(userID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	return nil
}
