package mycache

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewNodePool(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6")
	hash.Add("4")
	hash.Add("2")
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if peer, ok := hash.Search(k); !ok || peer != v {
			t.Errorf("Asking for %s, should have yielded %s but got %s", k, v, peer)
		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if peer, ok := hash.Search(k); !ok || peer != v {
			t.Errorf("Asking for %s, should have yielded %s but got %s", k, v, peer)
		}
	}

}
