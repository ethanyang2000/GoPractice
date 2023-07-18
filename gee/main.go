package main

import (
	gee "gee/mygee"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})
	r.Use(gee.Recovery())
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":8080")
}
