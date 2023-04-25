package mycache

import (
	"net/http"
	"strings"
)

type cacheServer struct {
	Prefix   string
	BasePath string
}

func NewServer(prefix string, basepath string) *cacheServer {
	return &cacheServer{
		Prefix:   prefix,
		BasePath: basepath,
	}
}

func (s *cacheServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
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
