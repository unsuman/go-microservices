package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/unsuman/go-microservices/aggregator/client"
	"github.com/unsuman/go-microservices/types"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type InvoiceHandler struct {
	c client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{c: c}
}

func main() {
	listenAddr := flag.String("listenaddr", ":9000", "HTTP gateway listen address")
	flag.Parse()

	c := client.NewGRPCClient("localhost:50051")
	invoiceHandler := NewInvoiceHandler(c)
	http.HandleFunc("/invoice", makeHandlerFunc(invoiceHandler.handleInvoice))
	http.Handle("/metrics", promhttp.Handler())
	logrus.Printf("HTTP gateway listening on %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func (i *InvoiceHandler) handleInvoice(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("invalid method")
	}
	obuIDStr := r.URL.Query().Get("obuid")
	obuID, err := strconv.ParseInt(obuIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid obuid: %v", err)
	}

	inv, err := i.c.GetInvoice(context.Background(), &types.InvoiceRequest{ObuID: obuID})
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
