package urlshortener

import (
	"errors"
	"log"
	"math/rand"
	"strings"
	"sync"
)

var ErrKeyGenerationFailed = errors.New("key_generation_failed")

func NewManager() *Manager {
	return &Manager{
		urlShortcuts: make(map[string]string),
	}
}

type Manager struct {
	mu           sync.RWMutex
	urlShortcuts map[string]string // short url key -> full url
}

func (m *Manager) CreateShortcut(fullURL string) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"
	const keyLength = 5
	const maxAttempts = 5

	for attempt := 0; attempt < maxAttempts; attempt++ {
		var keyBuilder strings.Builder
		for i := 0; i < keyLength; i++ {
			char := alphabet[rand.Intn(len(alphabet))]
			_ = keyBuilder.WriteByte(char) // string builder never fails writes
		}
		key := keyBuilder.String()

		succeeded := m.insertIfKeyIsNotUsed(key, fullURL)
		if !succeeded {
			log.Printf("Got collision for key %s. Retry key generation...", key)
			continue
		}
		return key, nil
	}
	return "", ErrKeyGenerationFailed
}

func (m *Manager) insertIfKeyIsNotUsed(key string, fullURL string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.urlShortcuts[key]; ok {
		return false
	}
	m.urlShortcuts[key] = fullURL
	return true
}

func (m *Manager) ResolveShortcut(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, found := m.urlShortcuts[key]
	return url, found
}
