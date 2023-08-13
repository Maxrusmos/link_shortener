package flagpkg

import (
	"sync"
	"testing"
)

func TestFlag(t *testing.T) {
	flag := NewFlag("initial")

	// Test initial value
	if value := flag.GetValue(); value != "initial" {
		t.Errorf("Expected initial value 'initial', but got '%s'", value)
	}

	// Test setting and getting value
	flag.SetValue("new value")
	if value := flag.GetValue(); value != "new value" {
		t.Errorf("Expected value 'new value', but got '%s'", value)
	}
}

func TestSharedFlagConcurrentAccess(t *testing.T) {
	flag := GetSharedFlag()

	var wg sync.WaitGroup
	const goroutines = 100

	// Concurrently set and get values from shared flag
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			flag.SetValue(string(rune(index)))
			value := flag.GetValue()
			if value != string(rune(index)) {
				t.Errorf("Expected value '%d', but got '%s'", index, value)
			}
		}(i)
	}

	wg.Wait()
}

func TestSharedFlagSetValue(t *testing.T) {
	flag := GetSharedFlag()

	// Test setting value in shared flag
	flag.SetValue("test value")
	if value := flag.GetValue(); value != "test value" {
		t.Errorf("Expected value 'test value', but got '%s'", value)
	}
}
