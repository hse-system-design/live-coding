package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

type handler struct {
	mu    sync.RWMutex
	alive bool
	ready bool
}

func (s *handler) HandleLivenessProbe(rw http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	alive := s.alive
	s.mu.RUnlock()

	if alive {
		rw.WriteHeader(http.StatusOK)
		return
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *handler) HandleReadinessProbe(rw http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	ready := s.ready
	s.mu.RUnlock()

	if ready {
		rw.WriteHeader(http.StatusOK)
		return
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *handler) HandleSetProbe(rw http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	if rawAlive := params.Get("alive"); rawAlive != "" {
		alive, err := strconv.ParseBool(rawAlive)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		s.Mutate(func() {
			s.alive = alive
		})
	}
	if rawReady := params.Get("ready"); rawReady != "" {
		ready, err := strconv.ParseBool(rawReady)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		s.Mutate(func() {
			s.ready = ready
		})
	}

	rw.WriteHeader(http.StatusOK)
}

func (s *handler) Mutate(mutator func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	mutator()
}

func main() {
	hndl := &handler{}
	go func() {
		time.Sleep(1 * time.Second)
		hndl.Mutate(func() {
			hndl.alive = true

		})
	}()
	go func() {
		time.Sleep(5 * time.Second)
		hndl.Mutate(func() {
			hndl.ready = true
		})
	}()

	router := http.NewServeMux()
	router.HandleFunc("/hello", func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("Hello from probe-server!\n"))
	})
	router.HandleFunc("/alive", hndl.HandleLivenessProbe)
	router.HandleFunc("/ready", hndl.HandleReadinessProbe)
	router.HandleFunc("/set-probe", hndl.HandleSetProbe)

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}
	_ = server.ListenAndServe()
}
