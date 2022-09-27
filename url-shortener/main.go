package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"url-shortener/grpcapi"
	"url-shortener/httpapi"
	"url-shortener/urlshortener"
	"url-shortener/urlshortener/inmemoryimpl"
	"url-shortener/urlshortener/mongoimpl"
)

const (
	modeInMemory = "in-memory"
	modeMongo    = "mongo"
)

var flagMode = flag.String("mode", modeInMemory, fmt.Sprintf("Storage mode. Possible values: %q, %q", modeInMemory, modeMongo))
var flagMongoAddr = flag.String("mongo-addr", "mongodb://localhost:27017", "Address of MongoDB to connect to")

func main() {
	flag.Parse()
	var manager urlshortener.Manager
	switch *flagMode {
	case modeInMemory:
		manager = inmemoryimpl.NewManager()
	case modeMongo:
		manager = mongoimpl.NewManager(*flagMongoAddr)
	default:
		log.Fatalf("Unexpected mode flag: %q", *flagMode)
	}

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
