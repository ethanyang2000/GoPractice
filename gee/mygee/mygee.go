package mygee

import (
	"net/http"
	"path"
	"strings"
	"text/template"
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
	htmlTemplate *template.Template
	htmlFuncMap  template.FuncMap
}

func New() *Engine {
	e := &Engine{r: NewRouter()}
	e.RouterGroup = &RouterGroup{
		engine:       e,
		basePath:     "",
		relativePath: "",
	}
	e.group = append(e.group, e.RouterGroup)
	return e
}

func (group *RouterGroup) Group(path string) *RouterGroup {
	g := &RouterGroup{
		engine:       group.engine,
		parent:       group,
		relativePath: path,
		basePath:     strings.Join([]string{group.basePath, path}, ""),
	}
	group.engine.group = append(group.engine.group, g)
	return g
}

func (gourp *RouterGroup) Use(middleware ...HandlerFunc) {
	gourp.middleware = append(gourp.middleware, middleware...)
}

func (g *RouterGroup) GET(str string, h HandlerFunc) {
	str = strings.Join([]string{g.basePath, str}, "")
	g.engine.r.addRoute("GET", str, h)
}

func (g *RouterGroup) POST(str string, h HandlerFunc) {
	str = strings.Join([]string{g.basePath, str}, "")
	g.engine.r.addRoute("POST", str, h)
}

func (g *RouterGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	c.engine = g.engine
	g.engine.r.handle(c, g)
}

func (g *RouterGroup) Run(port string) (err error) {
	return http.ListenAndServe(port, g.engine)
}

func (g *RouterGroup) Static(relativePath string, filePath string) {
	absolutePath := path.Join(g.basePath, relativePath)
	h := g.createStaticHandler(absolutePath, filePath)
	URLPath := path.Join(relativePath, "/*filepath")
	g.GET(URLPath, h)
}

func (g *RouterGroup) createStaticHandler(absolutePath string, filePath string) HandlerFunc {
	fileServer := http.StripPrefix(absolutePath, http.FileServer(http.Dir(filePath)))
	return func(c *Context) {
		file := c.Params["filepath"]
		if _, err := http.Dir(filePath).Open(file); err != nil {
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
