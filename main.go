package main

import (
	gee "mygee"
)

func main() {
	e := gee.New()
	e.GET("/ping", func(c *gee.Context) {
		name := c.Query("name")
		id := c.Query("id")
		c.JSON(200, gee.H{
			"name": name,
			"id":   id,
		})
	})
	e.Run(":8080")
}
