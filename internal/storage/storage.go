package storage

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"link_shortener/internal/dbwork"
	filework "link_shortener/internal/fileWork"
	"link_shortener/internal/shortenurl"
	"log"
	"os"
	"sync"
)

type URLStorage interface {
	AddURL(key string, url string) error
	GetURL(key string) (string, error)
	AddURLSH(url string) (string, error)
	Ping() error
}

type MapURLStorage struct {
	urls  map[string]string
	mutex sync.Mutex
}

func NewMapURLStorage() URLStorage {
	return &MapURLStorage{
		urls: make(map[string]string),
	}
}

func (s *MapURLStorage) AddURL(key string, url string) error {
	log.Println("MapURLStorageADDURL")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, found := s.urls[key]; found {
		return errors.New("key already exists")
	}
	s.urls[key] = url
	return nil
}

func (s *MapURLStorage) AddURLSH(url string) (string, error) {
	shortURL := shortenurl.Shortener(url)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.urls[shortURL] = url
	return shortURL, nil
}

func (s *MapURLStorage) GetURL(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	url, found := s.urls[key]
	if !found {
		return "", errors.New("key not found")
	}
	return url, nil
}

func (s *MapURLStorage) Ping() error {
	return errors.New("Ping is not supported for MapURLStorage")
}

type FileURLStorage struct {
	filePath string
	mutex    sync.Mutex
}

func NewFileURLStorage(filePath string) URLStorage {
	return &FileURLStorage{
		filePath: filePath,
	}
}

func (s *FileURLStorage) AddURL(key string, url string) error {
	log.Println("FileURLStorageADDURL")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	dataToWrite := filework.JSONURLs{
		ShortURL:  key,
		OriginURL: url,
	}
	data, err := json.Marshal(dataToWrite)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	_, err = file.WriteString(string(data) + "\n")
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func (s *FileURLStorage) AddURLSH(url string) (string, error) {
	shortURL := shortenurl.Shortener(url)
	log.Println("FileURLStorageADDURL")
	s.mutex.Lock()
	defer s.mutex.Unlock()

	urlToWrite := filework.JSONURLs{
		ShortURL:  shortURL,
		OriginURL: url,
	}

	err := filework.WriteURLsToFile(s.filePath, urlToWrite)
	if err != nil {
		log.Println("err is", err)
		return "", err
	}

	return shortURL, nil
}

type JSONURLs struct {
	ShortURL  string `json:"shortURL"`
	OriginURL string `json:"originURL"`
}

func (s *FileURLStorage) GetURL(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	file, err := os.OpenFile(s.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var data JSONURLs
	for scanner.Scan() {
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &data)
		if err != nil {
			return "", err
		}
		if data.ShortURL == key {
			return data.OriginURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return data.OriginURL, nil
}

func (s *FileURLStorage) Ping() error {
	return errors.New("Ping is not supported for FileURLStorage")
}

type DatabaseURLStorage struct {
	db    *sql.DB
	mutex sync.Mutex
}

func NewDatabaseURLStorage(db *sql.DB) URLStorage {
	return &DatabaseURLStorage{
		db: db,
	}
}

func (s *DatabaseURLStorage) AddURL(key string, url string) error {
	shortURL := shortenurl.Shortener(url)
	err := dbwork.AddURL(s.db, shortURL, url)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return nil
}

func (s *DatabaseURLStorage) AddURLSH(url string) (string, error) {
	shortURL := shortenurl.Shortener(url)
	err := dbwork.AddURL(s.db, shortURL, url)
	if err != nil {
		return "", err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return shortURL, nil
}

func (s *DatabaseURLStorage) GetURL(key string) (string, error) {
	originalURL, err := dbwork.GetOriginalURL(s.db, key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return originalURL, nil
}

func (s *DatabaseURLStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}
