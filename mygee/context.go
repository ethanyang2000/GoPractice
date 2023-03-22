package mygee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Method  string
	Path    string
	Params  map[string]string
	handler []HandlerFunc
	index   int
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: req,
		Method:  req.Method,
		Path:    req.URL.Path,
		index:   -1,
	}
}

func (c *Context) Next() {
	c.index++
	l := len(c.handler)
	for ; c.index < l; c.index++ {
		c.handler[c.index](c)
	}
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	} else {
		c.Status(code)
	}
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get((key))
}
