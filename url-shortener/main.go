package main

import (
	"log"
	"url-shortener/httpapi"
	"url-shortener/urlshortener"
)

func main() {
	manager := urlshortener.NewManager()
	srv := httpapi.NewServer(manager)
	log.Printf("Start serving on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
