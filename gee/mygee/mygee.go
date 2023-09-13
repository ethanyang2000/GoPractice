package mygee

import (
	"net/http"
	"path"
	"text/template"
)

type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	r      *router
	groups []*RouterGroup
}

type RouterGroup struct {
	engine       *Engine
	fullPath     string
	groupPath    string
	parent       *RouterGroup
	middleware   []HandlerFunc
	htmlTemplate *template.Template
	htmlFuncMap  template.FuncMap
}

func New() *Engine {
	e := &Engine{r: NewRouter()}
	e.RouterGroup = &RouterGroup{
		engine:    e,
		fullPath:  "",
		groupPath: "",
	}
	e.groups = append(e.groups, e.RouterGroup)
	return e
}

func (group *RouterGroup) Group(path string) *RouterGroup {
	g := &RouterGroup{
		engine:    group.engine,
		parent:    group,
		groupPath: path,
		fullPath:  group.fullPath + path,
	}
	group.engine.groups = append(group.engine.groups, g)
	return g
}

func (gourp *RouterGroup) Use(middleware ...HandlerFunc) {
	gourp.middleware = append(gourp.middleware, middleware...)
}

func (g *RouterGroup) GET(str string, h HandlerFunc) {
	g.engine.r.addRoute("GET", g.fullPath+str, h)
}

func (g *RouterGroup) POST(str string, h HandlerFunc) {
	g.engine.r.addRoute("POST", g.fullPath+str, h)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	c.engine = e
	e.r.handle(c, e)
}

func (e *Engine) Run(port string) (err error) {
	return http.ListenAndServe(port, e)
}

func (g *RouterGroup) Static(relativeRoute string, resourcePath string) {
	routePath := path.Join(g.fullPath, relativeRoute)
	h := g.createStaticHandler(routePath, resourcePath)
	URLPath := path.Join(routePath, "/*filepath")
	g.GET(URLPath, h)
}

func (g *RouterGroup) createStaticHandler(routePath string, resourcePath string) HandlerFunc {
	fileServer := http.StripPrefix(routePath, http.FileServer(http.Dir(resourcePath)))
	return func(c *Context) {
		file := c.Params["filepath"]
		if _, err := http.Dir(resourcePath).Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (g *RouterGroup) SetFuncMap(funcMap template.FuncMap) {
	g.htmlFuncMap = funcMap
}

func (g *RouterGroup) LoadHTMLGlob(pattern string) {
	g.htmlTemplate = template.Must(template.New("").Funcs(g.htmlFuncMap).ParseGlob(pattern))
}
