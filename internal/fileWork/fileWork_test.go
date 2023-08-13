package filework

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestWriteURLsToFile(t *testing.T) {
	filename := "test.txt"
	defer os.Remove(filename)

	dataToWrite := JSONURLs{
		ShortURL:  "short",
		OriginURL: "http://example.com",
	}

	err := WriteURLsToFile(filename, dataToWrite)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var writtenData JSONURLs
	err = json.Unmarshal(fileData, &writtenData)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if writtenData.ShortURL != dataToWrite.ShortURL {
		t.Errorf("Unexpected short URL: got %s, want %s", writtenData.ShortURL, dataToWrite.ShortURL)
	}
	if writtenData.OriginURL != dataToWrite.OriginURL {
		t.Errorf("Unexpected origin URL: got %s, want %s", writtenData.OriginURL, dataToWrite.OriginURL)
	}
}

func TestFindOriginURL(t *testing.T) {
	filename := "test.txt"
	defer os.Remove(filename)

	dataToWrite := JSONURLs{
		ShortURL:  "short",
		OriginURL: "http://example.com",
	}

	err := WriteURLsToFile(filename, dataToWrite)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	originURL, err := FindOriginURL(filename, "short")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if originURL != dataToWrite.OriginURL {
		t.Errorf("Unexpected origin URL: got %s, want %s", originURL, dataToWrite.OriginURL)
	}
}
