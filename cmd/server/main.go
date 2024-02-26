package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/docker/docker/client"

	"github.com/ilya-hontarau/distributed-object-storage/internal/gateway"
	"github.com/ilya-hontarau/distributed-object-storage/internal/server"
	"github.com/ilya-hontarau/distributed-object-storage/internal/storage"
	"github.com/ilya-hontarau/distributed-object-storage/internal/svcdiscovery"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer dockerClient.Close()
	dockerSvcDiscovery := svcdiscovery.NewDocker(dockerClient)
	configs, err := dockerSvcDiscovery.MinioConfigs(context.Background())
	if err != nil {
		panic(err)
	}
	storageNodes := make([]gateway.StorageNode, 0, len(configs))
	for _, cfg := range configs {
		logger.Info("config", slog.String("key", cfg.AccessKey), slog.String("secret", cfg.SecretKey))
		minio, err := storage.NewMinio(context.Background(), cfg.Addr, "default", cfg.AccessKey, cfg.SecretKey)
		if err != nil {
			panic(err)
		}
		storageNodes = append(storageNodes, minio)
	}
	gateway := gateway.NewGateway(storageNodes)
	handler := server.NewHTTPHandler(gateway, logger)
	// TODO: add graceful shutdown

	logger.Info("Starting service on port") // TODO: add port
	err = http.ListenAndServe(":3000", handler)
	if err != nil {
		panic(err)
	}
}
