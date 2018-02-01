package fakeapi

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Mux for fake api routes
func Mux(empty bool) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/primes", handler([]int{2, 3, 5, 7, 11, 13}, empty))
	m.HandleFunc("/fibo", handler([]int{1, 1, 2, 3, 5, 8, 13, 21}, empty))
	m.HandleFunc("/odd", handler([]int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23}, empty))
	m.HandleFunc("/rand", handler([]int{5, 17, 3, 19, 76, 24, 1, 5, 10, 34, 8, 27, 7}, empty))
	return m
}

func handler(numbers []int, empty bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		waitPeriod := rand.Intn(550)
		log.Printf("%s: waiting %dms.", r.URL.Path, waitPeriod)

		<-time.After(time.Duration(waitPeriod) * time.Millisecond)
		if empty {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{"numbers": numbers})
	}
}
