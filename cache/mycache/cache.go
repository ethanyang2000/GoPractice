package mycache

import (
	"sync"
)

type Cache struct {
	mu        sync.Mutex
	lru       *cache
	maxLen    int64
	onEvicted func(string, EntryValue)
}

func (c *Cache) Search(key string) (bv ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru != nil {
		v, ok := c.lru.Search(key)
		if v == nil {
			return ByteView{}, ok
		} else {
			return v.(ByteView), ok
		}
	} else {
		return
	}
}

func (c *Cache) Add(key string, value ByteView) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = NewCache(c.maxLen, c.onEvicted)
	}
	return c.lru.Add(key, value)
}

func (c *Cache) Delete(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.lru.Delete(key)
	return v.(ByteView), ok
}
