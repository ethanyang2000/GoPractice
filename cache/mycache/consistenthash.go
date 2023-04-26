package mycache

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func([]byte) uint32

type NodeMap struct {
	hash         Hash
	virtualNodes []uint32
	realNodes    map[uint32]string
	replicants   int
}

func NewNodePool(replicants int, hash Hash) *NodeMap {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &NodeMap{
		hash:         hash,
		replicants:   replicants,
		virtualNodes: make([]uint32, 0),
		realNodes:    make(map[uint32]string),
	}
}

func (nodes *NodeMap) Add(name string) {
	for i := 0; i < nodes.replicants; i++ {
		nodeName := strconv.Itoa(i) + name
		hashValue := uint32(nodes.hash([]byte(nodeName)))
		nodes.virtualNodes = append(nodes.virtualNodes, hashValue)
		nodes.realNodes[hashValue] = name
	}
	sort.Slice(nodes.virtualNodes, func(i, j int) bool {
		return nodes.virtualNodes[i] < nodes.virtualNodes[j]
	})

}

func (nodes *NodeMap) Search(key string) string {
	if len(nodes.realNodes) == 0 {
		return ""
	}
	hashValue := nodes.hash([]byte(key))
	idx := sort.Search(len(nodes.virtualNodes), func(i int) bool {
		return nodes.virtualNodes[i] >= hashValue
	})
	idx = idx % len(nodes.virtualNodes)
	return nodes.realNodes[nodes.virtualNodes[idx]]
}
