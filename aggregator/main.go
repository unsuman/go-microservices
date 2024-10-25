package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/unsuman/go-microservices/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3300", "HTTP server listen address")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLoggingMiddleware(svc)

	makeHTTPTransport(svc, *listenAddr)
}

func makeHTTPTransport(svc Aggregator, listenAddr string) {
	fmt.Println("HTTP transport listening on", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuid := r.URL.Query().Get("obuid")
		if obuid == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		obuidInt, err := strconv.ParseInt(obuid, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		distance, err := svc.CalculateInvoice(obuidInt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(distance)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
