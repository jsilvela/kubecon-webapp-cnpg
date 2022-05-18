package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/lib/pq"
)

const stockTableTpl string = `
<html>
<h3>Stonks</h3>
<h3>from most recently updated</h3>
<table>
{{ range . }}
	<tr>
		<td>{{ .Stock }}</td>
		<td>{{ .StockValue }}</td>
		<td>{{ .Date }}</td>
	</tr>
{{ end }}
</table>
</html>
`

const indexPage string = `
<html>
<h3>Hello KubeCon!</h3>

<ul>
	<li>Get <a href="/latest">the latest stock values</a></li>
	<li>Add <a href="/update">random stock values</a></li>
</ul>
</html>
`

// stockTicker represents the value of a stock at a given time
type stockTicker struct {
	Stock      string
	StockValue float64
	Date       time.Time
}

func main() {
	podName := os.Getenv("MY_POD_NAME")
	podEnv := os.Getenv("MY_POD_NAMESPACE")
	podIP := os.Getenv("MY_POD_IP")

	pgUser := os.Getenv("PG_USER")
	pgPass := os.Getenv("PG_PASSWORD")

	port := 8080

	var inside bool
	flag.BoolVar(&inside, "inside", false, "run webapp inside kind?")
	flag.Parse()

	var dbConnString string
	if inside {
		dbConnString = fmt.Sprintf("postgres://%s:%s@%s/app?sslmode=require",
			pgUser, pgPass, "cluster-example-rw")
	} else {
		dbConnString = fmt.Sprintf("postgres://%s:%s@%s/app?sslmode=require",
			pgUser, pgPass, "localhost")
	}

	stockTable, err := template.New("table").Parse(stockTableTpl)
	if err != nil {
		log.Fatalf("could not parse template: %v", err)
	}

	rand.Seed(time.Now().UnixNano())

	// HTTP ENDPOINTS

	// handle route using handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request", r.RequestURI, time.Now().UTC().Format(time.RFC3339))

		if r.Header.Get("Accept") == "application/json" {
			fmt.Fprintf(w, "Hello! pod: %s/%s ip: %s || %s\n%v",
				podName, podEnv, podIP, time.Now().Format(time.RFC3339), r.Header)
		} else {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, indexPage)
		}
	})

	http.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request", r.RequestURI, time.Now().UTC().Format(time.RFC3339))

		db, err := sql.Open("postgres", dbConnString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if errPing := db.Ping(); errPing != nil {
			http.Error(w, errPing.Error(), http.StatusInternalServerError)
			return
		}

		rows, err := db.Query(`
select stock, date, stock_value
from (
      select stock, rank() over wd as rank,
            first_value(date) over wd as date,
            first_value(stock_value) over wd as stock_value
      from stock_values
      window wd as (partition by stock order by date desc)
) as latest where rank = 1
order by date desc, stock;
`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var stocks []stockTicker
		for rows.Next() {
			var stockT stockTicker
			err = rows.Scan(&stockT.Stock, &stockT.Date, &stockT.StockValue)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			stocks = append(stocks, stockT)
		}

		if rErr := rows.Err(); rErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Header.Get("Accept") == "application/json" {
			jsonWriter := json.NewEncoder(w)
			err = jsonWriter.Encode(stocks)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "text/html")
			err = stockTable.Execute(w, stocks)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request", r.RequestURI, time.Now().UTC().Format(time.RFC3339))

		db, err := sql.Open("postgres", dbConnString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()

		if err := db.PingContext(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var success bool
		defer func() {
			if success {
				tx.Commit()
			} else {
				tx.Rollback()
			}
		}()

		for i := 0; i < 5; i++ {
			n := rand.Intn(50) + 1 // between 1 and 50, like our stocks
			stock := fmt.Sprintf("stock_%d", n)
			_, err = tx.Exec(
				`INSERT INTO stock_values(stock, date, stock_value) values ($1, $2, $3)`,
				stock, time.Now(), rand.Float64())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		success = true
		w.Write([]byte("done"))
	})

	// listen to port
	log.Println("ENV", podName, podEnv, podIP, pgPass, "cluster-example-rw", pgUser)
	log.Printf("running the server, listening on %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
