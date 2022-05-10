package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	podName := os.Getenv("MY_POD_NAME")
	podEnv := os.Getenv("MY_POD_NAMESPACE")
	podIP := os.Getenv("MY_POD_IP")
	// handle route using handler function
	count := 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! pod: %s/%s ip: %s || %d", podName, podEnv, podIP, count)
		count++
		log.Println("pew!!")
	})

	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	age := 21
	rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)
	_ = rows

	// listen to port
	log.Println("ENV", podName, podEnv, podIP)
	log.Println("running the server, listening on 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
