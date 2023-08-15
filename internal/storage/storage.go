package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"link_shortener/internal/dbwork"
	filework "link_shortener/internal/fileWork"
	"link_shortener/internal/shortenurl"
	"log"
	"sync"
)

type URLStorage interface {
	AddURL(key string, url string) error
	GetURL(key string) (string, error)
	Ping() error
}

type URLStorageWithAddURLSH interface {
	URLStorage
	AddURLSH(key string) (string, error)
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
	urlToWrite := filework.JSONURLs{
		ShortURL:  key,
		OriginURL: url,
	}
	err := filework.WriteURLsToFile(s.filePath, urlToWrite)
	if err != nil {
		log.Println("err is", err)
		return err
	}

	return nil
}

func (s *FileURLStorage) GetURL(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	originURL, err := filework.FindOriginURL(s.filePath, key)
	if err != nil {
		return "", err
	}

	return originURL, nil
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
