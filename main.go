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

	pgUser := os.Getenv("PG_USER")
	pgPass := os.Getenv("PG_PASSWORD")
	pgService := "cluster-example-rw"

	port := 8080
	// handle route using handler function
	count := 0
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! pod: %s/%s ip: %s || %d", podName, podEnv, podIP, count)
		count++
		log.Println("pew!!")
	})

	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		// postgres://$(USER):$(PASSWORD)@cluster-fast-failover-rw/app?sslmode=require&connect_timeout=2
		// connStr := "user=app dbname=app sslmode=verify-full"

		log.Println("ENV", podName, podEnv, podIP, "pass:", pgPass, pgService, "user:", pgUser)

		connStr := fmt.Sprintf("postgres://%s:%s@%s/app?sslmode=require", pgUser, pgPass, pgService)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errPing := db.Ping(); errPing != nil {
			http.Error(w, errPing.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := db.Query("SELECT bar, baz FROM foo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		fmt.Fprintf(w, "Hello! pod: %s/%s ip: %s || %d", podName, podEnv, podIP, count)

		for rows.Next() {
			var (
				bar int
				baz string
			)
			err = rows.Scan(&bar, &baz)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "row: %d -> %s", bar, baz)
		}

		if rErr := rows.Err(); rErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		count++
		log.Println("db pew!!")
	})

	// listen to port
	log.Println("ENV", podName, podEnv, podIP, pgPass, pgService, pgUser)
	log.Printf("running the server, listening on %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
