package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.Path("/metrics").Handler(promhttp.Handler())

	srv := &http.Server{Addr: "0.0.0.0:2112", Handler: r}
	log.Fatal(srv.ListenAndServe())
}
