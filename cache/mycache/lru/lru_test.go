package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int64 {
	return int64(len(d))
}

func TestSearchAndRemove(t *testing.T) {
	lru := NewCache(int64(10), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Search("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Search("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
	lru.Delete("key1")
	if _, ok := lru.Search("key1"); ok {
		t.Fatal("failed to delete cache")
	}
}

func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Search("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(e Entry) {
		keys = append(keys, e.Key)
	}
	lru := NewCache(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
