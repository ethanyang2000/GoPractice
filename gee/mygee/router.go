package mygee

import (
	"net/http"
	"strings"
)

type router struct {
	routerMap RouterMap
}

func NewRouter() *router {
	return &router{routerMap: make(RouterMap)}
}

func (r *router) addRoute(method string, str string, h HandlerFunc) bool {
	prefixes := ParsePrefix(str)
	if _, ok := r.routerMap[method]; !ok {
		r.routerMap[method] = NewRoot()
	}
	return r.routerMap[method].Insert(str, prefixes, 0, h)
}

func ParsePrefix(pattern string) []string {
	ptns := strings.Split(pattern, "/")
	prefixes := make([]string, 0)
	for idx, ptn := range ptns {
		if strings.HasPrefix(ptn, "*") {
			prefixes = append(prefixes, strings.Join(ptns[idx:], "/"))
			break
		}
		prefixes = append(prefixes, ptn)
	}
	return prefixes
}

func (r *router) findRoute(c *Context) HandlerFunc {
	prefixes := ParsePrefix(c.Path)
	if _, ok := r.routerMap[c.Method]; !ok {
		return func(c *Context) {
			c.String(http.StatusNotFound, "NOT FOUND")
		}
	}
	if node, ok := r.routerMap[c.Method].Search(prefixes, 0); !ok {
		return func(c *Context) {
			c.String(http.StatusNotFound, "NOT FOUND")
		}
	} else {
		params := make(map[string]string)
		patterns := strings.Split(node.Pattern, "/")
		for idx, pattern := range patterns {
			if strings.HasPrefix(pattern, ":") {
				params[string([]byte(pattern)[1:])] = prefixes[idx]
			}
			if strings.HasPrefix(pattern, "*") {
				patterns[idx] = string([]byte(patterns[idx])[1:])
				paramKey := strings.Join(patterns[idx:], "/")
				params[paramKey] = strings.Join(prefixes[idx:], "/")
			}
		}
		c.Params = params
		return node.handler
	}
}

func (r *router) handle(c *Context, g *RouterGroup) {
	handler := r.findRoute(c)
	for _, group := range g.engine.group {
		if strings.HasPrefix(c.Path, group.basePath) {
			c.handler = append(c.handler, group.middleware...)
		}
	}
	c.handler = append(c.handler, handler)
	c.Next()
}
