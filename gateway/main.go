package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenaddr", ":6000", "HTTP gateway listen address")
	flag.Parse()

	http.HandleFunc("/invoice", makeHandlerFunc(handleInvoice))
	fmt.Println("HTTP gateway listening on", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handleInvoice(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, map[string]string{"message": "invoice"})
}

func makeHandlerFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
