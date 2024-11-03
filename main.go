package main

import (
	"log"
	"net/http"

	server "server/src"
)

func main() {
	http.HandleFunc("/", server.HomePage)
	http.HandleFunc("/details/{id}", server.SecondPage)
	http.HandleFunc("/lastpage/{id}", server.LastPage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("Server running at http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
