package fileWork

import (
	"bufio"
	"encoding/json"
	"link_shortener/internal/storage"
	"os"
)

type JsonURLs struct {
	ShortURL  string `json:"shortURL"`
	OriginURL string `json:"originURL"`
}

func CheckIfURLExistsInFile(filename string, urls JsonURLs) (bool, error) {
	existingData, err := ReadJSONFile(filename)
	if err != nil {
		return false, err
	}
	for _, d := range existingData {
		if d.ShortURL == urls.ShortURL && d.OriginURL == urls.OriginURL {
			return true, nil
		}
	}
	return false, nil
}

func ReadJSONFile(filename string) ([]JsonURLs, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []JsonURLs
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var d JsonURLs
		err := json.Unmarshal(scanner.Bytes(), &d)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func WriteURLsToFile(filename string, dataToWrite JsonURLs) error {
	data, err := json.Marshal(dataToWrite)
	if err != nil {
		return err
	}

	isExist, err := CheckIfURLExistsInFile(filename, dataToWrite)
	if err != nil {
		return err
	}

	if isExist {
		return nil
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(string(data) + "\n")
	if err != nil {
		return err
	}
	return nil
}

func ReadDataFromFile(filename string, storage storage.URLStorage) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var data JsonURLs
		err := json.Unmarshal(scanner.Bytes(), &data)
		if err != nil {
			return err
		}
		err = storage.AddURL(data.ShortURL, data.OriginURL)
		if err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
