package main

import (
	"log"
	"time"

	gee "gee/mygee"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", 200, c.Request.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(func(c *gee.Context) {
		log.Printf("global middleware used")
	}) // global midlleware
	r.GET("/ping/:name", func(c *gee.Context) {
		c.JSON(200, c.Params)
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/geektutu
			c.JSON(200, c.Params)
		})
	}

	r.Run(":8080")
}
