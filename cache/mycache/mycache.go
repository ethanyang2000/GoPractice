package mycache

import (
	lru "cache/mycache/lru"
	"fmt"
	"sync"
)

type Group struct {
	MainCache *Cache
	Name      string
	peers     *HTTPPool
	loader    *CallGroup
	Getter    getter
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxLen int64, getter getter, onEvicted func(lru.Entry)) *Group {
	g := &Group{
		MainCache: &Cache{
			maxLen:    maxLen,
			onEvicted: onEvicted,
		},
		Name:   name,
		Getter: getter,
		loader: &CallGroup{},
	}
	mu.Lock()
	defer mu.Unlock()
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	return groups[name]
}

func (g *Group) Add(key string, value ByteView) {
	g.MainCache.Add(key, value)
}

func (g *Group) Delete(key string) (ByteView, bool) {
	return g.MainCache.Delete(key)
}

func (g *Group) Search(key string) (ByteView, error) {
	if bv, ok := g.MainCache.Search(key); ok {
		return bv, nil
	} else {
		return g.load(key)
	}
}

func (g *Group) load(key string) (ByteView, error) {
	data, err := g.loader.Call(key, func() (interface{}, error) {
		if g.peers != nil {
			if data, err := g.peers.Search(g.Name, key); err == nil {
				return ByteView{b: data}, err
			}
		}
		return g.getFromLocal(key)
	})
	if err != nil {
		return ByteView{}, err
	}
	return data.(ByteView), err
}

func (g *Group) getFromLocal(key string) (ByteView, error) {
	v, err := g.Getter.Get(key)
	if err != nil {
		return ByteView{}, fmt.Errorf("get data from local database failed")
	}
	bv := ByteView{b: v}
	g.MainCache.Add(key, bv)
	return bv, nil
}

func (g *Group) RegisterPeers(peers *HTTPPool) {
	if g.peers == nil {
		g.peers = peers
	} else {
		panic("duplicant peers initiated")
	}
}
