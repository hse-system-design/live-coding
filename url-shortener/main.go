package main

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"url-shortener/grpcapi"
	"url-shortener/httpapi"
	"url-shortener/ratelimit"
	"url-shortener/urlshortener"
	"url-shortener/urlshortener/inmemoryimpl"
	"url-shortener/urlshortener/mongoimpl"
	rediscached "url-shortener/urlshortener/rediscachedimpl"
)

const (
	modeInMemory = "in-memory"
	modeMongo    = "mongo"
	modeCached   = "cached"
)

var flagMode = flag.String("mode", modeInMemory, fmt.Sprintf("Storage mode. Possible values: %q, %q", modeInMemory, modeMongo))
var flagMongoAddr = flag.String("mongo-addr", "mongodb://localhost:27017", "Address of MongoDB to connect to")
var flagRedisAddr = flag.String("redis-addr", "127.0.0.1:6379", "Address of Redis to connect to")

func main() {
	flag.Parse()

	redisClient := redis.NewClient(&redis.Options{Addr: *flagRedisAddr})
	limiterFactory := ratelimit.NewFactory(redisClient)

	var manager urlshortener.Manager
	switch *flagMode {
	case modeInMemory:
		manager = inmemoryimpl.NewManager()
	case modeMongo:
		manager = mongoimpl.NewManager(*flagMongoAddr)
	case modeCached:
		manager = rediscached.NewManager(
			redisClient,
			mongoimpl.NewManager(*flagMongoAddr),
		)
	default:
		log.Fatalf("Unexpected mode flag: %q", *flagMode)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		srv := httpapi.NewServer(manager, limiterFactory)
		log.Printf("Start serving HTTP at %s", srv.Addr)
		log.Fatal(srv.ListenAndServe())
	}()
	go func() {
		defer wg.Done()
		grpcapi.RunGRPCServer(manager)
	}()

	wg.Wait()
}
