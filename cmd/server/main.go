package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/docker/docker/client"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/sync/errgroup"

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

	logger.Info("Starting server...", slog.Int("port", envCfg.Port))
	err = runServer(envCfg, handler)
	if err != nil {
		panic(err)
	}
}

func runServer(cfg Config, handler http.Handler) error {
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: time.Second,
	}
	exitCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var group errgroup.Group
	group.Go(func() error {
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to listen and serve: %w", err)
		}
		return nil
	})
	group.Go(func() error {
		<-exitCtx.Done()
		ctx, cancelF := context.WithTimeout(context.Background(), time.Second)
		defer cancelF()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return nil
	})
	err := group.Wait()
	if err != nil {
		return fmt.Errorf("failed to wait for group: %w", err)
	}
	return nil
}
