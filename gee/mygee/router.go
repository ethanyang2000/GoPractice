package mygee

import (
	"net/http"
	"strings"
)

type RouterMap map[string]*node
type router struct {
	routerMap RouterMap
}

func NewRouter() *router {
	return &router{routerMap: make(RouterMap)}
}

func (r *router) addRoute(method string, path string, h HandlerFunc) error {
	if _, ok := r.routerMap[method]; !ok {
		r.routerMap[method] = NewRoot()
	}
	return r.routerMap[method].Insert(strings.Split(path, "/"), h)
}

func (r *router) findRoute(c *Context) HandlerFunc {
	if _, ok := r.routerMap[c.Method]; !ok {
		return func(c *Context) {
			c.String(http.StatusNotFound, " NOT FOUND")
		}
	}
	results := []string{}
	params := []string{}
	if node, err := r.routerMap[c.Method].Search(strings.Split(c.Path, "/"), &results, &params); err != nil {
		return func(c *Context) {
			c.String(http.StatusNotFound, "NOT FOUND")
		}
	} else {
		p := make(map[string]string)
		idx := 0
		for _, r := range results {
			if strings.HasPrefix(r, "*") || strings.HasPrefix(r, ":") {
				p[string([]byte(r)[1:])] = params[idx]
				idx += 1
			}
		}
		c.Params = p
		return node.handler
	}
}

func (r *router) handle(c *Context, e *Engine) {
	handler := r.findRoute(c)
	for _, group := range e.groups {
		if strings.HasPrefix(c.Path, group.fullPath) {
			c.handlers = append(c.handlers, group.middleware...)
		}
	}
	c.handlers = append(c.handlers, handler)
	c.Next()
}
