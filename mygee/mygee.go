package mygee

import(
	"net/http"
	"fmt"
)

type HandlerFunc func(w http.ResponseWriter, req *http.Request)
type RouterMap map[string]map[string]HandlerFunc

type engine struct{
	routeMap RouterMap
}

func New() *engine{
	return &engine{routeMap: make(RouterMap)}
}

func (e *engine) addRoute(method string, str string, h HandlerFunc){
	if _,ok := e.routeMap[method];!ok{
		e.routeMap[method] = make(map[string]HandlerFunc)
	}
	e.routeMap[method][str] = h
}

func  (e *engine) GET(str string, h HandlerFunc){
	e.addRoute("GET", str, h)
}

func (e *engine) POST(str string, h HandlerFunc){
	e.addRoute("POST", str, h)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	if _, ok := e.routeMap[req.Method];!ok{
		fmt.Fprintf(w, "HTTP ERROR: 404 NOT FOUND FOR %s \n", req.URL.Path)
		return
	}
	if h, ok := e.routeMap[req.Method][req.URL.Path];!ok{
		fmt.Fprintf(w, "HTTP ERROR: 404 NOT FOUND FOR %s \n", req.URL.Path)
	}else{
		h(w, req)
	}
}

func (e *engine) Run(port string) (err error){
	return http.ListenAndServe(port, e)
}
