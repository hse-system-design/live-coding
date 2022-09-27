package inmemoryimpl

import (
	"context"
	"log"
	"sync"
	"url-shortener/urlshortener"
)

func NewManager() *Manager {
	return &Manager{
		urlShortcuts: make(map[string]string),
	}
}

type Manager struct {
	mu           sync.RWMutex
	urlShortcuts map[string]string // short url key -> full url
}

func (m *Manager) CreateShortcut(_ context.Context, fullURL string) (string, error) {
	const maxAttempts = 5

	for attempt := 0; attempt < maxAttempts; attempt++ {
		key := urlshortener.GenerateKey()

		succeeded := m.insertIfKeyIsNotUsed(key, fullURL)
		if !succeeded {
			log.Printf("Got collision for key %s. Retry key generation...", key)
			continue
		}
		return key, nil
	}
	return "", urlshortener.ErrKeyGenerationFailed
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

func (m *Manager) ResolveShortcut(_ context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, found := m.urlShortcuts[key]
	if !found {
		return "", urlshortener.ErrNotFound
	}
	return url, nil
}
