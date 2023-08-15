package shortenurl

import (
	"testing"
)

func TestShortener(t *testing.T) {
	tests := []struct {
		originalURL string
		expected    string
	}{
		{"http://example.com", "a9b9f043"},
		{"https://google.com", "99999ebc"},
		{"http://stackoverflow.com", "57f4dad4"},
	}

	for _, test := range tests {
		result := Shortener(test.originalURL)
		if result != test.expected {
			t.Errorf("Unexpected result for originalURL '%s': got %s, want %s", test.originalURL, result, test.expected)
		}
	}
}
