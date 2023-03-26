package mycache

import "sync"

type Group struct {
	MainCache *Cache
	Name      string
	Getter    getter
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxLen int64, getter getter, onEvicted func(string, EntryValue)) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		MainCache: &Cache{
			maxLen:    maxLen,
			onEvicted: onEvicted,
			lru:       NewCache(maxLen, onEvicted),
		},
		Name:   name,
		Getter: getter,
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	return groups[name]
}

func (g *Group) Add(key string, value ByteView) bool {
	return g.MainCache.Add(key, value)
}

func (g *Group) Delete(key string) (ByteView, bool) {
	return g.MainCache.Delete(key)
}

func (g *Group) Search(key string) (ByteView, bool) {
	if bv, ok := g.MainCache.Search(key); ok {
		return bv, ok
	} else {
		return g.load(key)
	}
}

func (g *Group) load(key string) (ByteView, bool) {
	v, ok := g.Getter.Get(key)
	bv := ByteView{b: v.Slice()}
	g.MainCache.Add(key, bv)
	return bv, ok
}
