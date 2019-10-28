package controllers

import (
	"Demo/if/framework"
	"Demo/if/models"
	"fmt"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	framework.GetIcecream().DB.LogMode(true)
	blog := models.Blog{}
	blog2 := models.Blog{}
	framework.GetIcecream().DB.First(&blog)
	framework.GetIcecream().DB.Last(&blog2)
	blog.Content = ""
	blog2.Content = ""
	fmt.Printf("blog: %+v\n", blog)
	fmt.Printf("blog2: %+v\n", blog2)
	blog.ID = 0
	framework.GetIcecream().DB.Create(&blog)

	data := map[string]interface{}{
		"title": "index",
		"Name":  "Bob",
		"Age":   18,
		"Books": []string{"books1", "books2", "books3"},
		"Classmates": map[string]string{
			"lily": "16",
			"bill": "17",
			"jobs": "18",
		},
	}

	t, _ := template.ParseFiles("views/index.html", "views/header.html")
	t.Execute(w, data)
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "statis/favicon.ico")
}
