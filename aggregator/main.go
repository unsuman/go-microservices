package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/unsuman/go-microservices/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("listenaddr", ":3300", "HTTP server listen address")
	grpcListenAddr := flag.String("grpcaddr", ":3301", "GRPC server listen address")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLoggingMiddleware(svc)

	go makeHTTPTransport(svc, *httpListenAddr)
	makeGRPCServer(svc, *grpcListenAddr)
}

func makeHTTPTransport(svc Aggregator, listenAddr string) {
	fmt.Println("HTTP transport listening on", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func makeGRPCServer(svc Aggregator, listenAddr string) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}

	server := grpc.NewServer()
	fmt.Println("GRPC transport listening on", listenAddr)

	types.RegisterAggregatorServer(server, NewGRPCAggregator(svc))

	server.Serve(ln)
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
