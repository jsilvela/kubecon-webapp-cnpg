package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// handle route using handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to new server!")
		log.Println("pew!!")
	})

	// listen to port
	log.Println("running the server, listening on 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
