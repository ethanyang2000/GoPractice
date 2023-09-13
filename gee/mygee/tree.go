package mygee

import (
	"errors"
	"strings"
)

type node struct {
	isPath   bool
	nodePath string
	Children []*node
	handler  HandlerFunc
}

func NewRoot() *node {
	return &node{isPath: false, nodePath: "", Children: make([]*node, 0)}
}

func (n *node) Search(targets []string, results *[]string, params *[]string) (*node, error) {
	if strings.HasPrefix(n.nodePath, "*") {
		*params = append(*params, strings.Join(targets, "/"))
		return n, nil
	}
	if len(targets) == 0 {
		if n.isPath {
			return n, nil
		} else {
			return nil, errors.New("route do not exist")
		}
	}
	for _, child := range n.Children {
		if child.nodePath == targets[0] {
			*results = append(*results, n.nodePath)
			return n.Search(targets[1:], results, params)
		} else {
			if strings.HasPrefix(child.nodePath, "*") {
				*results = append(*results, n.nodePath)
				return n.Search(targets[1:], results, params)
			}
			if strings.HasPrefix(child.nodePath, ":") {
				*results = append(*results, n.nodePath)
				*params = append(*params, targets[0])
				return n.Search(targets[1:], results, params)
			}
		}
	}
	return nil, errors.New("route do not exist")
}

func (n *node) Insert(targets []string, handler HandlerFunc) error {
	// register a dynamic route ":" with other routes in the same path is allowed,
	// but the results are unpredicable
	if len(targets) == 0 {
		n.isPath = true
		n.handler = handler
		return nil
	}

	if len(n.Children) == 1 && (strings.HasPrefix(n.Children[0].nodePath, "*")) {
		if n.Children[0].nodePath == targets[0] {
			n.Children[0].isPath = true
			n.Children[0].handler = handler
			return nil
		} else {
			return errors.New("cannot add new routes to existing dynamic routes \"*\"")
		}
	}

	if len(n.Children) > 0 && (strings.HasPrefix(targets[0], "*")) {
		return errors.New("cannot add dynamic routes \"*\" to existing routes")
	}

	for _, child := range n.Children {
		if targets[0] == child.nodePath {
			return child.Insert(targets[1:], handler)
		}
	}
	new_node := &node{nodePath: targets[0], Children: make([]*node, 0), isPath: false}
	n.Children = append(n.Children, new_node)
	return new_node.Insert(targets[1:], handler)
}
