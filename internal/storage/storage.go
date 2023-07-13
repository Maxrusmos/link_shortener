package storage

import (
	"errors"
	"link_shortener/internal/shortenurl"
)

type URLStorage interface {
	AddURL(key string, url string) error
	GetURL(key string) (string, error)
	AddURLSH(url string) (string, error)
}

type mapURLStorage struct {
	urls map[string]string
}

func NewMapURLStorage() URLStorage {
	return &mapURLStorage{
		urls: make(map[string]string),
	}
}

func (s *mapURLStorage) AddURLSH(url string) (string, error) {
	shortURL := shortenurl.Shortener(url)
	s.urls[shortURL] = url
	return shortURL, nil
}

func (s *mapURLStorage) AddURL(key string, url string) error {
	if _, found := s.urls[key]; found {
		return errors.New("key already exists")
	}
	s.urls[key] = url
	return nil
}

func (s *mapURLStorage) GetURL(key string) (string, error) {
	url, found := s.urls[key]
	if !found {
		return "", errors.New("key not found")
	}
	return url, nil
}
