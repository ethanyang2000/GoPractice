package mygee

import (
	"strings"
)

type node struct {
	Pattern  string
	Prefix   string
	Children []*node
	handler  HandlerFunc
}

func NewRoot() *node {
	return &node{Pattern: "", Prefix: "/", Children: make([]*node, 0)}
}

func (n *node) Search(prefixes []string, depth int) (*node, bool) {
	if depth == len(prefixes) || strings.HasPrefix(n.Prefix, "*") {
		if n.Pattern == "" {
			return nil, false
		}
		return n, true
	}

	for _, child := range n.Children {
		if prefixes[depth] == child.Prefix || strings.HasPrefix(child.Prefix, "*") || strings.HasPrefix(child.Prefix, ":") {
			if ans, ok := child.Search(prefixes, depth+1); ok {
				return ans, true
			}
		}
	}
	return nil, false
}

func (n *node) Insert(pattern string, prefixes []string, depth int, handler HandlerFunc) bool {
	if len(prefixes) == depth {
		n.Pattern = pattern
		n.handler = handler
		return true
	}

	for _, child := range n.Children {
		if prefixes[depth] == child.Prefix {
			child.Insert(pattern, prefixes, depth+1, handler)
		} else if (strings.HasPrefix(child.Prefix, ":") && strings.HasPrefix(prefixes[depth], ":")) || (strings.HasPrefix(child.Prefix, "*") && strings.HasPrefix(prefixes[depth], "*")) {
			child.Prefix = prefixes[depth]
			child.Insert(pattern, prefixes, depth+1, handler)
			//simply replace the original dynamic route when a new one is set
		}
	}
	new_node := &node{Prefix: prefixes[depth], Children: make([]*node, 0), Pattern: ""}
	n.Children = append(n.Children, new_node)
	new_node.Insert(pattern, prefixes, depth+1, handler)

	return false
}
