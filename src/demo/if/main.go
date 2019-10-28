package main

import (
	"Demo/if/models"
	_ "Demo/if/routers"
	"fmt"
	"net/http"
	"reflect"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	addr string `json:"_"`
}

func main() {
	blog := models.NewBlog()
	blog.GetBlogs()
	fmt.Printf("blogs: %+v\n", blog)
	t := reflect.TypeOf(*blog)
	fmt.Printf("blogT: %+v\n", t)
	blog.SetAll(blog, map[string]interface{}{"id": 16, "title": "xxxxx"})
	//blog.Save(blog, map[string]interface{}{"id": 16, "title": "xxxxx"})
	fmt.Printf("blogs: %+v\n", blog)
	//title := blog.Get(blog, "Title")
	//fmt.Printf("blog title: %+v\n", title)

	server := &http.Server{Addr: "127.0.0.1:8888"}
	http.HandleFunc("/favicon.ico", favicon)
	server.ListenAndServe()
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "statics/favicon.ico")
}
