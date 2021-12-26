package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

var buckets = map[int]*rate.Limiter{}

var maxConcurrency int

func main() {
	if concurrency, ok := os.LookupEnv("CONCURRENCY"); ok {
		val, err := strconv.Atoi(concurrency)
		if err != nil {
			maxConcurrency = 1
			log.Println("Could not parse \"" + concurrency + "\" as concurrency! Continuing with maxConcurrency = 1")
		} else {
			maxConcurrency = val
		}
	} else {
		maxConcurrency = 1
		log.Println("No CONCURRENCY found in environment! Continuing with maxConcurrency = 1")
	}

	if maxConcurrency == 0 {
		log.Fatalln("CONCURRENCY may not be zero!")
	}

	for i := 0; i < maxConcurrency; i++ {
		buckets[i] = rate.NewLimiter(rate.Every(time.Second*6), 1)
	}

	r := http.NewServeMux()

	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		shard := 0
		if val := r.URL.Query().Get("shard"); val != "" {
			v2, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Warn: could not convert shard ID to int. Continuing as shard 0")
			} else {
				shard = v2
			}
		} else {
			log.Println("Warn: shard ID not found in query. Continuing as shard 0")
		}

		bucket := shard % maxConcurrency

		log.Println("Waiting for identify on shard " + strconv.Itoa(shard))

		// Wait(context.Context) blocks the current goroutine until the next identify slot
		// and cancels the request if the HTTP request dies prematurely
		if err := buckets[bucket].Wait(r.Context()); err == nil {
			log.Println("shard " + strconv.Itoa(shard) + " can connect now")
			rw.Write([]byte("You are free to connect now! :)"))
		} else {
			log.Println("Shard " + strconv.Itoa(shard) + " Wait error: " + err.Error())
		}
	})

	addr := ":8080"

	if val, ok := os.LookupEnv("ADDR"); ok {
		addr = val
	}

	log.Println("Listening on", addr)
	log.Fatalln(http.ListenAndServe(addr, r))
}
