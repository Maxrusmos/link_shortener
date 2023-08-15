package filework

import (
	"bufio"
	"encoding/json"
	"os"
)

type JSONURLs struct {
	ShortURL  string `json:"shortURL"`
	OriginURL string `json:"originURL"`
}

func WriteURLsToFile(filename string, dataToWrite JSONURLs) error {
	data, err := json.Marshal(dataToWrite)
	if err != nil {
		return err
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

func FindOriginURL(filename string, shortURL string) (string, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var data JSONURLs
		line := scanner.Text()                     // Читаем строку файла
		err := json.Unmarshal([]byte(line), &data) // Разбираем JSON из строки
		if err != nil {
			return "", err
		}
		if data.ShortURL == shortURL {
			return data.OriginURL, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", err
}
