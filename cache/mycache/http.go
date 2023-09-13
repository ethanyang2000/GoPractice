package mycache

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HTTPPool struct {
	selfPath string
	basePath string
	peers    *PeerManager
}

func NewPool(selfPath string, basepath string, replicants int, hashFunc Hash) *HTTPPool {
	return &HTTPPool{
		selfPath: selfPath,
		basePath: basepath,
		peers:    NewPeerManager(replicants, hashFunc),
	}
}

func (s *HTTPPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL.Path

	if !(strings.HasPrefix(url, s.basePath)) {
		http.Error(w, url+" Not Found", http.StatusNotFound)
	}

	urls := strings.SplitN(string([]byte(url)[len(s.basePath)+1:]), "/", 2)

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

func (s *HTTPPool) Add(peer, url string) error {
	return s.peers.Add(peer, url)
}

func (s *HTTPPool) Search(group, key string) ([]byte, error) {
	u, err := s.peers.Search(key)

	if err != nil || u == s.selfPath {
		return []byte{}, errors.New("failed to get cache from peers")
	}
	url := u + fmt.Sprintf(
		"%v/%v/%v",
		s.basePath,
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
