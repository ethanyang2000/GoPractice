package main

import (
	gee "mygee"
)

func main() {
	e := gee.New()
	group := e.Group("g1")
	group2 := group.Group("g2")
	e.GET("ping/:idx/*id", func(c *gee.Context) {
		c.JSON(200, c.Params)
	})
	group2.GET("pong/:idx/*ids", func(c *gee.Context) {
		c.JSON(200, c.Params)
	})
	e.Run(":8080")
}
