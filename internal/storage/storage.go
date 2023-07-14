package storage

import (
	"errors"
	"link_shortener/internal/shortenurl"
	"sync"
)

type URLStorage interface {
	AddURL(key string, url string) error
	GetURL(key string) (string, error)
	AddURLSH(url string) (string, error)
}

type mapURLStorage struct {
	urls  map[string]string
	mutex sync.Mutex
}

func NewMapURLStorage() URLStorage {
	return &mapURLStorage{
		urls: make(map[string]string),
	}
}

func (s *mapURLStorage) AddURLSH(url string) (string, error) {
	shortURL := shortenurl.Shortener(url)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.urls[shortURL] = url
	return shortURL, nil
}

func (s *mapURLStorage) AddURL(key string, url string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, found := s.urls[key]; found {
		return errors.New("key already exists")
	}
	s.urls[key] = url
	return nil
}

func (s *mapURLStorage) GetURL(key string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	url, found := s.urls[key]
	if !found {
		return "", errors.New("key not found")
	}
	return url, nil
}
