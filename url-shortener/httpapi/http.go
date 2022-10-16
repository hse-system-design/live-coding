package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
	"url-shortener/ratelimit"
	"url-shortener/urlshortener"
)

func NewHTTPHandler(
	manager urlshortener.Manager,
	limiterFactory *ratelimit.Factory,
) *HTTPHandler {

	return &HTTPHandler{
		manager:      manager,
		createLimit:  limiterFactory.NewLimiter("post_url", 10*time.Second, 2),
		resolveLimit: limiterFactory.NewLimiter("get_url", 1*time.Minute, 10),
	}
}

type HTTPHandler struct {
	manager urlshortener.Manager

	createLimit  *ratelimit.Limiter
	resolveLimit *ratelimit.Limiter
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
	if isRateLimited(h.createLimit, rw, r) {
		return
	}

	var data CreateShortcutRequest

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	key, err := h.manager.CreateShortcut(r.Context(), data.Url)
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

func (h *HTTPHandler) ResolveURL(rw http.ResponseWriter, r *http.Request) {
	if isRateLimited(h.resolveLimit, rw, r) {
		return
	}

	key := strings.Trim(r.URL.Path, "/")

	url, err := h.manager.ResolveShortcut(r.Context(), key)
	if errors.Is(err, urlshortener.ErrNotFound) {
		http.NotFound(rw, r)
	}
	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusPermanentRedirect)
}

func NewServer(
	manager urlshortener.Manager,
	rateLimitFactory *ratelimit.Factory,
) *http.Server {
	r := mux.NewRouter()

	handler := NewHTTPHandler(manager, rateLimitFactory)

	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)
	r.PathPrefix("/").Methods(http.MethodOptions).HandlerFunc(corsPreflightHandler)
	r.HandleFunc("/{shortUrl:\\w{5}}", handler.ResolveURL).Methods(http.MethodGet)
	r.HandleFunc("/api/urls", handler.CreateShortcut).Methods(http.MethodPost)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv
}

type responseWriter struct {
	http.ResponseWriter
	Status int
}

func (rw *responseWriter) WriteHeader(status int) {
	if rw.Status == 0 {
		rw.Status = status
	}
	rw.ResponseWriter.WriteHeader(status)
}

func isRateLimited(limiter *ratelimit.Limiter, rw http.ResponseWriter, r *http.Request) bool {
	canDo, err := limiter.CanDoAt(r.Context(), time.Now())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return true
	}
	if !canDo {
		http.Error(rw, "rate limit exceeded", http.StatusTooManyRequests)
		return true
	}
	return false
}

func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		wrapper := &responseWriter{ResponseWriter: rw}

		start := time.Now()
		h.ServeHTTP(wrapper, r)
		elapsed := time.Now().Sub(start)

		log.Printf("%s %s: %d %s", r.Method, r.URL, wrapper.Status, elapsed)
	})
}

func corsPreflightHandler(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
	rw.Header().Set("Access-Control-Expose-Headers", "*")
	rw.WriteHeader(http.StatusNoContent)
}

func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			h.ServeHTTP(rw, r)
			return
		}

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "*")
		rw.Header().Set("Access-Control-Expose-Headers", "*")
		h.ServeHTTP(rw, r)
	})
}
