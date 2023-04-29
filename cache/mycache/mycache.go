package mycache

import (
	lru "cache/mycache/lru"
	"fmt"
	"sync"
)

type Group struct {
	MainCache *Cache
	Name      string
	Getter    getter
	peers     PeerPicker
	loader    *CallGroup
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, maxLen int64, getter getter, onEvicted func(string, lru.EntryValue)) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		MainCache: &Cache{
			maxLen:    maxLen,
			onEvicted: onEvicted,
			lru:       lru.NewCache(maxLen, onEvicted),
		},
		Name:   name,
		Getter: getter,
		loader: &CallGroup{},
	}
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
			if peer, ok := g.peers.Pick(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
			}
		}
		return g.getFromLocal(key)
	})
	if err != nil {
		return ByteView{}, err
	}
	return data.(ByteView), err
}

func (g *Group) getFromPeer(getter PeerGetter, key string) (ByteView, error) {
	value, err := getter.Search(g.Name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{B: value}, nil
}

func (g *Group) getFromLocal(key string) (ByteView, error) {
	v, err := g.Getter.Get(key)
	if err != nil {
		return ByteView{}, fmt.Errorf("get data from local database failed")
	}
	newV := make([]byte, len(v))
	copy(newV, v)
	bv := ByteView{B: newV}
	g.MainCache.Add(key, bv)
	return bv, nil
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers == nil {
		g.peers = peers
	} else {
		panic("duplicant peers initiated")
	}
}
