package mygee

import (
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)
type RouterMap map[string]*node

type Engine struct {
	*RouterGroup
	r     *router
	group []*RouterGroup
}

type RouterGroup struct {
	engine       *Engine
	basePath     string
	relativePath string
	parent       *RouterGroup
	middleware   []HandlerFunc
}

func New() *Engine {
	e := &Engine{r: NewRouter()}
	e.RouterGroup = &RouterGroup{
		engine:       e,
		basePath:     "",
		relativePath: "",
	}
	return e
}

func (group *RouterGroup) Group(path string) *RouterGroup {
	g := &RouterGroup{
		engine:       group.engine,
		parent:       group,
		relativePath: path,
		basePath:     strings.Join([]string{group.basePath, path}, "/"),
	}
	group.engine.group = append(group.engine.group, g)
	return g
}

func (g *RouterGroup) GET(str string, h HandlerFunc) {
	str = strings.Join([]string{g.basePath, str}, "/")
	g.engine.r.addRoute("GET", str, h)
}

func (g *RouterGroup) POST(str string, h HandlerFunc) {
	str = strings.Join([]string{g.basePath, str}, "/")
	g.engine.r.addRoute("POST", str, h)
}

func (g *RouterGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	g.engine.r.handle(c)
}

func (g *RouterGroup) Run(port string) (err error) {
	return http.ListenAndServe(port, g.engine)
}
