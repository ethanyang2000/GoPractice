package main

import(
	"net/http"
	gee "mygee"
	"fmt"
)

func main(){
	e := gee.New()
	e.GET("/ping", func(w http.ResponseWriter, req *http.Request){
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	e.POST("/ping", func(w http.ResponseWriter, req *http.Request){
		fmt.Fprintf(w, "POST URL.Path = %q\n", req.URL.Path)
	})
	e.Run(":8080")
}