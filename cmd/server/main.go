package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/ilya-hontarau/distributed-object-storage/internal/gateway"
	"github.com/ilya-hontarau/distributed-object-storage/internal/server"
	"github.com/ilya-hontarau/distributed-object-storage/internal/storage"
)

func main() {
	minio, err := storage.NewMinio(context.Background(), "localhost:9000", "default", "ring", "treepotato")
	if err != nil {
		panic(err)
	}
	gateway := gateway.NewGateway([]gateway.StorageNode{minio})
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	handler := server.NewHTTPHandler(gateway, logger)
	err = http.ListenAndServe(":3000", handler)
	if err != nil {
		panic(err)
	}
}
