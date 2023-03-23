package mycache

import "container/list"

type cache struct {
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

func NewCache(maxLength int64, onEvicted func(string, EntryValue)) (c *cache) {
	return &cache{
		MaxLength: maxLength,
		Length:    0,
		CacheMap:  make(map[string]*list.Element),
		List:      list.New(),
		onEvicted: onEvicted,
	}
}

func (c *cache) Add(key string, value EntryValue) ([]*entry, bool) {
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
		return c.RemoveOld()
	}
	return nil, true
}

func (c *cache) Delete(key string) (*entry, bool) {
	if value, ok := c.CacheMap[key]; ok {
		c.List.Remove(value)
		delete(c.CacheMap, key)
		c.onEvicted(key, value.Value.(*entry).Value)
		c.Length -= int64(len(key)) + int64(value.Value.(*entry).Value.Len())
		return value.Value.(*entry), true
	} else {
		return nil, false
	}
}

func (c *cache) removeOld() (*entry, bool) {
	r := c.List.Back()
	kv := r.Value.(*entry)
	delete(c.CacheMap, kv.Key)
	c.List.Remove(r)
	c.onEvicted(kv.Key, kv.Value)
	c.Length -= int64(len(kv.Key)) + int64(kv.Value.Len())
	return kv, true
}

func (c *cache) RemoveOld() ([]*entry, bool) {
	removed := make([]*entry, 0)
	for c.Length > c.MaxLength {
		if r, ok := c.removeOld(); ok {
			removed = append(removed, r)
		} else {
			return nil, false
		}
	}
	return removed, true
}

func (c *cache) Get(key string) (*entry, bool) {
	if value, ok := c.CacheMap[key]; ok {
		c.List.MoveToFront(value)
		return value.Value.(*entry), true
	} else {
		return nil, false
	}
}
