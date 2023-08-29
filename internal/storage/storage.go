package storage

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	config "link_shortener/internal/configs"
	"link_shortener/internal/dbwork"
	"link_shortener/internal/shortenurl"
	"log"
	"os"
	"sync"
)

type URLStorage interface {
	AddURL(key string, url string, userID string) error
	GetURL(key string) (string, error)
	AddURLSH(url string, userID string) (string, error)
	GetOriginalURL(key string) (string, bool)
	Ping() error
	GetAllURLs(userID string) ([]map[string]string, error)
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

func (s *MapURLStorage) GetOriginalURL(key string) (string, bool) {
	return key, false
}

func (s *MapURLStorage) AddURL(key string, url string, userID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, found := s.urls[key]; found {
		fmt.Println("key already exists")
	}
	s.urls[key] = url
	return nil
}

func (s *MapURLStorage) AddURLSH(url string, userID string) (string, error) {
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

func (s *MapURLStorage) GetAllURLs(userID string) ([]map[string]string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	urls := make([]map[string]string, 0)
	for shortURL, originalURL := range s.urls {
		url := make(map[string]string)
		url["original_url"] = originalURL
		url["short_url"] = conf.BaseURL + "/" + shortURL
		urls = append(urls, url)
	}
	return urls, nil
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

func (s *FileURLStorage) GetOriginalURL(key string) (string, bool) {
	return key, false
}

type JSONURLs struct {
	ShortURL  string `json:"shortURL"`
	OriginURL string `json:"originURL"`
	UserID    string `json:"userID"`
}

var conf = config.GetConfig()

func (s *FileURLStorage) AddURL(key string, url string, userID string) error {
	log.Println(s.filePath, key, url)
	s.mutex.Lock()
	dataToWrite := JSONURLs{
		ShortURL:  conf.BaseURL + "/" + key,
		OriginURL: url,
		UserID:    userID,
	}
	data, err := json.Marshal(dataToWrite)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(conf.FileStore, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(string(data) + "\n")
	defer s.mutex.Unlock()
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func (s *FileURLStorage) AddURLSH(url string, userID string) (string, error) {
	shortURL := shortenurl.Shortener(url)
	log.Println("FileURLStorageADDURL")
	s.mutex.Lock()
	dataToWrite := JSONURLs{
		ShortURL:  conf.BaseURL + "/" + shortURL,
		OriginURL: url,
		UserID:    userID,
	}
	data, err := json.Marshal(dataToWrite)
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.OpenFile(conf.FileStore, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	_, err = file.WriteString(string(data) + "\n")
	if err != nil {
		fmt.Println(err)
	}
	defer s.mutex.Unlock()

	if err != nil {
		log.Println("err is", err)
		return "", err
	}

	return shortURL, nil
}

func (s *FileURLStorage) GetURL(key string) (string, error) {
	log.Println(s.filePath, key)
	s.mutex.Lock()
	file, err := os.OpenFile(conf.FileStore, os.O_RDONLY|os.O_CREATE, 0666)
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
			log.Println("originalURL SHORT", data.ShortURL)
			log.Println("originalURL LONG", data.OriginURL)
			return data.OriginURL, nil
		}
	}
	defer s.mutex.Unlock()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return data.OriginURL, nil
}

func (s *FileURLStorage) GetAllURLs(userID string) ([]map[string]string, error) {
	urls := []JSONURLs{}

	file, err := os.OpenFile(conf.FileStore, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var url JSONURLs
		err := json.Unmarshal(scanner.Bytes(), &url)
		if err != nil {
			continue
		}

		if url.UserID == userID {
			urls = append(urls, url)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var urlMaps []map[string]string
	for _, url := range urls {
		urlMap := make(map[string]string)
		urlMap["short_url"] = url.ShortURL
		urlMap["original_url"] = url.OriginURL
		urlMaps = append(urlMaps, urlMap)
	}

	return urlMaps, nil
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

func (s *DatabaseURLStorage) GetOriginalURL(key string) (string, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	originalURL, err := dbwork.GetOriginalURL(s.db, key)
	if err == sql.ErrNoRows {
		return "такой записи не существует", false
	}
	if originalURL != "" {
		return originalURL, true
	}
	return originalURL, false
}

func (s *DatabaseURLStorage) AddURL(key string, url string, userID string) error {
	shortURL := shortenurl.Shortener(url)
	err := dbwork.AddURL(s.db, conf.BaseURL+"/"+shortURL, url, userID)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return nil
}

func (s *DatabaseURLStorage) AddURLSH(url string, userID string) (string, error) {
	shortURL := shortenurl.Shortener(url)
	err := dbwork.AddURL(s.db, shortURL, url, userID)
	if err != nil {
		return "", err
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return shortURL, nil
}

func (s *DatabaseURLStorage) GetURL(key string) (string, error) {
	originalURL, err := dbwork.GetOriginalURL(s.db, conf.BaseURL+"/"+key)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return originalURL, nil
}

type URL struct {
	shortURL    string
	originalURL string
}

func (s *DatabaseURLStorage) GetAllURLs(userID string) ([]map[string]string, error) {
	query := "SELECT short_url, original_url FROM shortened_urls WHERE user_id = $1"
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := []URL{}
	for rows.Next() {
		var url URL
		err := rows.Scan(&url.shortURL, &url.originalURL)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	var urlMaps []map[string]string
	for _, url := range urls {
		urlMap := make(map[string]string)
		urlMap["short_url"] = url.shortURL
		urlMap["original_url"] = url.originalURL
		urlMaps = append(urlMaps, urlMap)
	}

	return urlMaps, nil
}

func (s *DatabaseURLStorage) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}
