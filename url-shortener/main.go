package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var ErrKeyGenerationFailed = errors.New("key_generation_failed")

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		urlShortcuts: make(map[string]string),
	}
}

type HTTPHandler struct {
	mu           sync.RWMutex
	urlShortcuts map[string]string // short url key -> full url
}

type CreateShortcutRequest struct {
	Url string `json:"url"`
}

type CreateShortcutResponse struct {
	Key string `json:"key"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *HTTPHandler) CreateShortcut(rw http.ResponseWriter, r *http.Request) {
	var data CreateShortcutRequest

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.generateKeyAndStoreURL(data.Url)
	var status int
	var response interface{}
	if err != nil {
		status = http.StatusInternalServerError
		response = ErrorResponse{
			Error: err.Error(),
		}
	} else {
		status = http.StatusOK
		response = CreateShortcutResponse{
			Key: key,
		}
	}
	rawResponse, _ := json.Marshal(response)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, _ = rw.Write(rawResponse)
}

func (h *HTTPHandler) generateKeyAndStoreURL(fullURL string) (string, error) {
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

		keyAlreadyUsed := false
		func() {
			h.mu.Lock()
			defer h.mu.Unlock()

			if _, ok := h.urlShortcuts[key]; ok {
				keyAlreadyUsed = true
			} else {
				h.urlShortcuts[key] = fullURL
			}
		}()
		if keyAlreadyUsed {
			log.Printf("Got collision for key %s. Retry key generation...", key)
			continue
		}
		return key, nil
	}
	return "", ErrKeyGenerationFailed
}

func (h *HTTPHandler) ResolveURL(rw http.ResponseWriter, r *http.Request) {
	key := strings.Trim(r.URL.Path, "/")

	url, found := h.getURLByKey(key)
	if !found {
		http.NotFound(rw, r)
		return
	}
	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusPermanentRedirect)
}

func (h *HTTPHandler) getURLByKey(key string) (string, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	url, found := h.urlShortcuts[key]
	return url, found
}

func main() {
	r := mux.NewRouter()

	handler := NewHTTPHandler()

	r.HandleFunc("/{shortUrl:\\w{5}}", handler.ResolveURL).Methods(http.MethodGet)
	r.HandleFunc("/api/urls", handler.CreateShortcut).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
