package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var pool *sql.DB

func main() {
	var err error
	pool, err = sql.Open("postgres", "postgres://postgres:password@localhost:5432/pg_isolation?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	err = pool.PingContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("databse connected...")

	http.HandleFunc("/", index)
	http.HandleFunc("/decrement", Decrement)
	http.HandleFunc("/concurrent", HandleConcurrent)

	errChan := make(chan error, 1)

	go func() {
		log.Println("listening on port 4000")
		err := http.ListenAndServe(":4000", nil)
		if err != nil {
			errChan <- err
		}
	}()

	log.Fatal(<-errChan)
}

func Query() {
	type Row struct {
		id       int64
		name     string
		quantity int64
	}

	var r Row

	var sql = `SELECT id, name, quantity FROM product`

	err := pool.QueryRow(sql).Scan(&r.id, &r.name, &r.quantity)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(r)
}
