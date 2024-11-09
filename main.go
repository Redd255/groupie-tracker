package main

import (
	"log"
	"net/http"

	server "server/src"
)

func main() {
	http.HandleFunc("/", server.HomePage)
	http.HandleFunc("/details/", server.SecondPage)
	http.HandleFunc("/lastpage/", server.LastPage)
	// this will be applicated in the next project "bonus"
	// http.Handle("/static/", http.FileServer(http.Dir("static")))
	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
