package mycache

import (
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("scores", 2<<10, GetFunc(
		func(key string) (ByteView, bool) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return ByteView{b: []byte(v)}, true
			}
			return ByteView{}, false
		}), nil)

	for k, v := range db {
		if view, ok := gee.Search(k); ok != true || view.String() != v {
			t.Fatal("failed to get value of Tom")
		} // load from callback function
		if _, ok := gee.Search(k); ok != true || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, ok := gee.Search("unknown"); ok == true {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
