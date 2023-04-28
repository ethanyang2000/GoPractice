package mycache

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HTTPPool struct {
	Prefix      string
	BasePath    string
	mu          sync.Mutex
	peers       *NodeMap
	httpGetters map[string]*httpGetter
}

func NewPool(prefix string, basepath string, replicants int, hashFunc Hash) *HTTPPool {
	return &HTTPPool{
		Prefix:      prefix,
		BasePath:    basepath,
		peers:       NewNodePool(replicants, hashFunc),
		httpGetters: make(map[string]*httpGetter),
	}
}

func (s *HTTPPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Path

	if !(strings.HasPrefix(url, s.BasePath)) {
		http.Error(w, url+" Not Found", http.StatusNotFound)
	}

	urls := strings.SplitN(string([]byte(url)[len(s.BasePath):]), "/", 2)

	if len(urls) != 2 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	groupName, Key := urls[0], urls[1]

	group := GetGroup(groupName)

	if group == nil {
		http.Error(w, "Cache Group Not Found", http.StatusNotFound)
		return
	}
	v, err := group.Search(Key)
	if err == nil {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(v.Slice())
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *HTTPPool) Pick(key string) (PeerGetter, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if peer := s.peers.Search(key); peer != "" && peer != s.Prefix {
		return s.httpGetters[peer], true
	}
	return nil, false
}

func (s *HTTPPool) Add(peers ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, peer := range peers {
		s.peers.Add(peer)
		s.httpGetters[peer] = &httpGetter{
			BasePath: peer + s.BasePath,
		}
	}
}

type httpGetter struct {
	BasePath string
}

func (g *httpGetter) Search(group, key string) ([]byte, error) {
	url := fmt.Sprintf(
		"%v%v/%v",
		g.BasePath,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("server return: %v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return bytes, err
}
