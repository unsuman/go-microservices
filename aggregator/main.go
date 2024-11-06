package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/unsuman/go-microservices/aggregator/client"
	"github.com/unsuman/go-microservices/types"
	"google.golang.org/grpc"
)

func main() {
	// httpListenAddr := flag.String("listenaddr", ":3300", "HTTP server listen address")
	grpcListenAddr := flag.String("grpcaddr", "localhost:50051", "GRPC server listen address")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLoggingMiddleware(svc)

	makeGRPCServer(svc, *grpcListenAddr)
	time.Sleep(time.Second * 5)
	// makeHTTPTransport(svc, *httpListenAddr)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := client.NewGRPCClient(*grpcListenAddr)
	if _, err := c.AggregatorClient.Aggregate(ctx, &types.AggregateRequest{
		ObuID: 1,
		Value: 10.0,
		Unix:  time.Now().Unix(),
	}); err != nil {
		fmt.Println("failed to call grpc client")
		return
	}
}

func makeHTTPTransport(svc Aggregator, listenAddr string) {
	fmt.Println("HTTP transport listening on", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func makeGRPCServer(svc Aggregator, listenAddr string) {
	ln, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		fmt.Println("failed to listen:", err)
		return
	}

	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	fmt.Println("GRPC transport listening on", listenAddr)

	types.RegisterAggregatorServer(server, NewGRPCAggregator(svc))

	if err = server.Serve(ln); err != nil {
		fmt.Println("failed to serve:", err)
		return
	}

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
