package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "myapp_http_duration_seconds",
	}, []string{"path"})
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		defer timer.ObserveDuration()

		next.ServeHTTP(rw, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.Use(prometheusMiddleware)

	r.Path("/metrics").Handler(promhttp.Handler())
	r.Path("/generate-pair").Methods(http.MethodPost).HandlerFunc(handleMineKey)

	srv := &http.Server{Addr: "0.0.0.0:2112", Handler: r}
	log.Fatal(srv.ListenAndServe())
}
