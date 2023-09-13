package mycache

import (
	lru "cache/mycache/lru"
	"sync"
)

type Cache struct {
	mu        sync.Mutex
	lru       *lru.LruCache
	maxLen    int64
	onEvicted func(lru.Entry)
}

func (c *Cache) Search(key string) (bv ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru != nil {
		v, ok := c.lru.Search(key)
		if !ok {
			return bv, ok
		} else {
			return v.(ByteView), ok
		}
	} else {
		return bv, false
	}
}

func (c *Cache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.maxLen, c.onEvicted)
	}
	c.lru.Add(key, value)
}

func (c *Cache) Delete(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.lru.Delete(key)
	return v.(ByteView), ok
}
