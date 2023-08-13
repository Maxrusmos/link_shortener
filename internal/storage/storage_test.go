package storage

import (
	"testing"
)

func TestNewMapURLStorage(t *testing.T) {
	storage := NewMapURLStorage()

	if storage == nil {
		t.Error("Expected non-nil storage")
	}

	_, err := storage.GetURL("key")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}

	err = storage.AddURL("key", "http://example.com")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	url, err := storage.GetURL("key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if url != "http://example.com" {
		t.Errorf("Unexpected URL: got %s, want %s", url, "http://example.com")
	}
}
