package main

import (
	"log"
	"sync"
	"url-shortener/grpcapi"
	"url-shortener/httpapi"
	"url-shortener/urlshortener/inmemoryimpl"
)

func main() {
	manager := inmemoryimpl.NewManager()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		srv := httpapi.NewServer(manager)
		log.Printf("Start serving HTTP at %s", srv.Addr)
		log.Fatal(srv.ListenAndServe())
	}()
	go func() {
		defer wg.Done()
		grpcapi.RunGRPCServer(manager)
	}()

	wg.Wait()
}
