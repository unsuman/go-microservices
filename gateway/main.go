package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/unsuman/go-microservices/aggregator/client"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type InvoiceHandler struct {
	c client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{c: c}
}

func main() {
	listenAddr := flag.String("listenaddr", ":6000", "HTTP gateway listen address")
	flag.Parse()

	c := client.NewGRPCClient("localhost:50051")
	invoiceHandler := NewInvoiceHandler(c)
	http.HandleFunc("/invoice", makeHandlerFunc(invoiceHandler.handleInvoice))
	fmt.Println("HTTP gateway listening on", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func (i *InvoiceHandler) handleInvoice(w http.ResponseWriter, r *http.Request) error {
	inv, err := i.c.GetInvoice(context.Background(), 1)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
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
