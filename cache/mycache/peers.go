package mycache

import (
	"errors"
	"sync"
)

type PeerManager struct {
	mu    sync.Mutex
	peers *NodeMap
	urls  map[string]string
}

func NewPeerManager(replicants int, hash Hash) *PeerManager {
	return &PeerManager{
		peers: NewNodePool(replicants, hash),
		urls:  make(map[string]string),
	}
}

func (p *PeerManager) Add(peer, url string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if url == "" {
		url = peer
	}
	if _, ok := p.urls[peer]; ok {
		return errors.New("peer already added")
	}
	p.peers.Add(peer)
	p.urls[peer] = url
	return nil
}

func (p *PeerManager) Search(key string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if res, ok := p.peers.Search(key); ok {
		return p.urls[res], nil
	} else {
		return "", errors.New("failed to find peer")
	}
}
