package storage

import (
	"testing"
)

func TestMapURLStorage(t *testing.T) {
	testStorage(t, NewMapURLStorage())
}

func TestFileURLStorage(t *testing.T) {
	testStorage(t, NewFileURLStorage("test_storage.json"))
}

func testStorage(t *testing.T, storage URLStorage) {
	storage.AddURL("abc123", "http://example.com", "user1")
	url, err := storage.GetURL("abc123")
	if err != nil {
		t.Errorf("GetURL returned an error: %v", err)
	}
	if url != "http://example.com" {
		t.Errorf("Expected URL to be %s, but got %s", "http://example.com", url)
	}

}
