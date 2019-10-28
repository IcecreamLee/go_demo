package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	server := &http.Server{Addr: "127.0.0.1:8080"}
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/tmpl", tmpl)
	http.HandleFunc("/json", json)
	server.ListenAndServe()
}

func index(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:     "c1",
		Value:    "Manning Publications Co",
		HttpOnly: true,
	}
	c2 := http.Cookie{
		Name:     "c2",
		Value:    "Manning Publications Co",
		HttpOnly: true,
	}
	w.Header().Set("Set-Cookie", c1.String())
	w.Header().Add("Set-Cookie", c2.String())

	fmt.Fprintf(w, "Hello World!!")
	//w.Write([]byte("Hello World!!"))
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login Page")
}

func json(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, `{"hello":"world"}`)
}

func tmpl(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl.html", "header.html")
	books := []string{"books1", "books2", "books3"}
	classmates := map[string]string{
		"lily": "16",
		"bill": "17",
		"jobs": "18",
	}
	user := map[string]interface{}{
		"Name":       "Bob",
		"Age":        18,
		"Books":      books,
		"Classmates": classmates,
	}
	t.Execute(w, user)
}
