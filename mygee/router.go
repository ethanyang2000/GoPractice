package mygee

import "net/http"

type router struct {
	routerMap RouterMap
}

func NewRouter() *router {
	return &router{routerMap: make(RouterMap)}
}

func (r *router) addRoute(method string, str string, h HandlerFunc) {
	if _, ok := r.routerMap[method]; !ok {
		r.routerMap[method] = make(map[string]HandlerFunc)
	}
	r.routerMap[method][str] = h
}

func (r *router) handle(c *Context) {
	if _, ok := r.routerMap[c.Method]; !ok {
		http.Error(c.Writer, "NOT FOUND", 404)
		return
	}
	if handler, ok := r.routerMap[c.Method][c.Path]; !ok {
		http.Error(c.Writer, "NOT FOUND", 404)
	} else {
		handler(c)
	}
}
