package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		defer wg.Done()
		go makeARequest()
	}

	wg.Wait()
}

func makeARequest() {
	res, err := http.Get("http://localhost:4000/decrement")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("response:", string(resByte))
}
