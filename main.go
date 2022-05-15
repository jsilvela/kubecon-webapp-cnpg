package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
		log.Println("request", r.RequestURI, time.Now().UTC().Format(time.RFC3339))
	})

	http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
		// postgres://$(USER):$(PASSWORD)@cluster-fast-failover-rw/app?sslmode=require&connect_timeout=2
		// connStr := "user=app dbname=app sslmode=verify-full"

		log.Println("request", r.RequestURI, time.Now().UTC().Format(time.RFC3339))

		// connStr := fmt.Sprintf("postgres://%s:%s@%s/app?sslmode=require", pgUser, pgPass, pgService)
		connStr := fmt.Sprintf("postgres://%s:%s@%s/app?sslmode=require", pgUser, pgPass, "localhost")
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errPing := db.Ping(); errPing != nil {
			http.Error(w, errPing.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := db.Query(`
select bond, date, factor
from (
      select bond, rank() over wd as rank,
            first_value(date) over wd as date,
            first_value(factor) over wd as factor
      from factors
      window wd as (partition by bond order by date desc)
) as latest where rank = 1;
`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		fmt.Fprintf(w, "Hello! pod: %s/%s ip: %s || %d", podName, podEnv, podIP, count)

		for rows.Next() {
			var (
				factor float64
				bond   string
				date   time.Time
			)
			err = rows.Scan(&bond, &date, &factor)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "row: %s: %e (%s)", bond, factor, date.Format(time.RFC3339))
		}

		if rErr := rows.Err(); rErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		count++
	})

	// listen to port
	log.Println("ENV", podName, podEnv, podIP, pgPass, pgService, pgUser)
	log.Printf("running the server, listening on %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
