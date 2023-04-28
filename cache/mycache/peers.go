package mycache

type PeerPicker interface {
	Pick(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	Search(group, key string) ([]byte, error)
}
