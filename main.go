package main

import (
	gee "mygee"
)

func main() {
	e := gee.New()
	e.GET("/ping/:idx/*id", func(c *gee.Context) {
		c.JSON(200, c.Params)
	})
	e.GET("/ping/:idx/*ids/fjaosjf/fasjo", func(c *gee.Context) {
		c.JSON(200, c.Params)
	})
	e.Run(":8080")
}
