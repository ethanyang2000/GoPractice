package main

import (
	"fmt"
	"mycache"
)

type Strings string

func (s Strings) Len() int64 {
	return int64(len(s))
}

func main() {
	k1 := "key1"
	k2 := "key2"
	v1 := "1234"
	v2 := "2345"
	ans := k1 + k2 + v1 + v2
	lru := mycache.NewCache(int64(len(ans)), func(s string, e mycache.EntryValue) {
		fmt.Println("callback called during deletion")
	})
	lru.Add("key1", Strings("1234"))
	lru.Add("key2", Strings("3456"))
	lru.Add("key3", Strings("456"))
	v, ok := lru.Get("key1")
	fmt.Println(v, ok)
}
