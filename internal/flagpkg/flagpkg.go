package flagpkg

import (
	"sync"
)

type Flag struct {
	value string
	mu    sync.RWMutex
}

func NewFlag(initialValue string) *Flag {
	return &Flag{value: initialValue}
}

func (f *Flag) SetValue(newValue string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.value = newValue
}

func (f *Flag) GetValue() string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.value
}

var sharedFlag = NewFlag("noF")

func GetSharedFlag() *Flag {
	return sharedFlag
}
