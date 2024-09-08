package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello, World!")
	fs := http.FileServer(http.Dir("static/css"))
	http.Handle("/", http.FileServer(http.Dir("htmx/")))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
