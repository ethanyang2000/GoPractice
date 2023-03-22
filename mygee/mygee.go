package mygee

import (
	"net/http"
)

type HandlerFunc func(c *Context)
type RouterMap map[string]*node

type Engine struct {
	r *router
	c *Context
}

func New() *Engine {
	return &Engine{r: NewRouter(), c: new(Context)}
}

func (e *Engine) GET(str string, h HandlerFunc) {
	e.r.addRoute("GET", str, h)
}

func (e *Engine) POST(str string, h HandlerFunc) {
	e.r.addRoute("POST", str, h)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	e.r.handle(c)
}

func (e *Engine) Run(port string) (err error) {
	return http.ListenAndServe(port, e)
}
