package main

import (
	mycache "cache/mycache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	mycache.NewGroup("scores", 2<<10, mycache.GetFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return []byte{}, fmt.Errorf("key not found in local database")
		}), nil)

	addr := "localhost:6666"
	peers := mycache.NewServer(addr, "/mycache/")
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
