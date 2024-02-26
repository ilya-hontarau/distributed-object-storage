package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/docker/docker/client"
	"github.com/kelseyhightower/envconfig"

	"github.com/ilya-hontarau/distributed-object-storage/internal/gateway"
	"github.com/ilya-hontarau/distributed-object-storage/internal/server"
	"github.com/ilya-hontarau/distributed-object-storage/internal/storage"
	"github.com/ilya-hontarau/distributed-object-storage/internal/svcdiscovery"
)

type Config struct {
	Port       int    `envconfig:"PORT" default:"3000"`
	BucketName string `envconfig:"BUCKET_NAME" default:"default"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var envCfg Config
	err := envconfig.Process("", &envCfg)
	if err != nil {
		panic(err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer dockerClient.Close()

	dockerSvcDiscovery := svcdiscovery.NewDocker(dockerClient)
	configs, err := dockerSvcDiscovery.MinioConfigs(context.TODO())
	if err != nil {
		panic(err)
	}
	storageNodes := make([]gateway.StorageNode, 0, len(configs))
	for _, cfg := range configs {
		minio, err := storage.NewMinio(context.TODO(), cfg.Addr, envCfg.BucketName, cfg.AccessKey, cfg.SecretKey)
		if err != nil {
			panic(err)
		}
		storageNodes = append(storageNodes, minio)
	}
	gateway := gateway.NewGateway(storageNodes)
	handler := server.NewHTTPHandler(gateway, logger)

	// TODO: add graceful shutdown
	logger.Info("Starting server...", slog.Int("port", envCfg.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", envCfg.Port), handler)
	if err != nil {
		panic(err)
	}
}
