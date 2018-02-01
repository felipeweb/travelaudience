package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/felipeweb/travelaudience/numbers"
)

func main() {
	listenAddr := flag.String("http.addr", ":8080", "http listen address")
	flag.Parse()
	http.HandleFunc("/numbers", numbers.OrderHandler)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
