package main

import (
	"log"
	"net/http"
	"sync"
)

func HandleConcurrent(w http.ResponseWriter, r *http.Request) {
	var before int64
	err := pool.QueryRow("SELECT quantity FROM product WHERE id = 1").Scan(&before)
	if err != nil {
		panic(err)
	}
	log.Printf("before update is %d", before)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			_, err := pool.Exec("UPDATE product SET quantity = quantity - 1 WHERE id = 1")
			if err != nil {
				panic(err)
			}
		}()
	}

	wg.Wait()

	var after int64
	err = pool.QueryRow("SELECT quantity FROM product WHERE id = 1").Scan(&after)
	if err != nil {
		panic(err)
	}
	log.Printf("after update is %d", after)

	w.Write([]byte("ok"))
}
