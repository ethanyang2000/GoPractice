package main

import (
	"cache/mycache"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var db = map[string]string{
	"1": "630",
	"2": "589",
	"3": "567",
}

func createGroup() *mycache.Group {
	return mycache.NewGroup("scores", 2<<10, mycache.GetFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), nil)
}

func startCacheServer(addr string, addrs []string, gee *mycache.Group) {
	peers := mycache.NewPool(addr, "/geecache", 50, func(b []byte) uint32 {
		i, _ := strconv.Atoi(string(b))
		return uint32(i)
	})
	for idx, a := range addrs {
		peers.Add(strconv.Itoa(idx), a)
	}
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *mycache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Search(key)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.Slice())

		}))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:5555"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}
