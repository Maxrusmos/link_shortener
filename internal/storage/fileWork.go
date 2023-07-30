package storage

import (
	"encoding/json"
	"io/ioutil"
)

type jsonURLs struct {
	shortURL  string
	originURL string
}

func WriteURLsToFile(filename string, dataToWrite jsonURLs) error {
	data, err := json.Marshal(dataToWrite)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadURLsFromFile(filename string) ([]mapURLStorage, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var urls []mapURLStorage
	err = json.Unmarshal(data, &urls)
	if err != nil {
		return nil, err
	}
	return urls, nil
}
