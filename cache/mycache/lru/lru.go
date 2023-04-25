package lru

import "container/list"

type LruCache struct {
	MaxLength int64
	Length    int64
	CacheMap  map[string]*list.Element
	List      *list.List
	onEvicted func(string, EntryValue)
}

type entry struct {
	Key   string
	Value EntryValue
}

type EntryValue interface {
	Len() int64
}

func NewCache(maxLength int64, onEvicted func(string, EntryValue)) (c *LruCache) {
	return &LruCache{
		MaxLength: maxLength,
		Length:    0,
		CacheMap:  make(map[string]*list.Element),
		List:      list.New(),
		onEvicted: onEvicted,
	}
}

func (c *LruCache) Add(key string, value EntryValue) {
	if valueOld, ok := c.CacheMap[key]; ok {
		kv := valueOld.Value.(*entry)
		c.List.MoveToFront(valueOld)
		c.Length += int64(value.Len()) - int64(kv.Value.Len())
		kv.Value = value
	} else {
		new_entry := &entry{Key: key, Value: value}
		c.List.PushFront(new_entry)
		c.CacheMap[key] = c.List.Front()
		c.Length += int64(value.Len()) + int64(len(key))
	}
	if c.Length > c.MaxLength {
		c.RemoveOld()
	}
}

func (c *LruCache) Delete(key string) (EntryValue, bool) {
	if ele, ok := c.CacheMap[key]; ok {
		c.List.Remove(ele)
		delete(c.CacheMap, key)
		if c.onEvicted != nil {
			c.onEvicted(key, ele.Value.(*entry).Value)
		}
		c.Length -= int64(len(key)) + int64(ele.Value.(*entry).Value.Len())
		return ele.Value.(*entry).Value, true
	} else {
		return nil, false
	}
}

func (c *LruCache) RemoveOld() {
	for c.Length > c.MaxLength {
		r := c.List.Back()
		if r != nil {
			kv := r.Value.(*entry)
			delete(c.CacheMap, kv.Key)
			c.List.Remove(r)
			if c.onEvicted != nil {
				c.onEvicted(kv.Key, kv.Value)
			}
			c.Length -= int64(len(kv.Key)) + int64(kv.Value.Len())
		}
	}
}

func (c *LruCache) Search(key string) (ev EntryValue, ok bool) {
	value, ok := c.CacheMap[key]
	if ok {
		c.List.MoveToFront(value)
		return value.Value.(*entry).Value, true
	} else {
		return nil, false
	}
}

func (c *LruCache) Len() int {
	return len(c.CacheMap)
}
