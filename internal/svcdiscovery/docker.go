package svcdiscovery

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const imageName = "minio/minio"

type Docker struct {
	client    *client.Client
	imageName string
}

func NewDocker(client *client.Client) *Docker {
	return &Docker{client: client, imageName: imageName}
}

type MinioConfig struct {
	AccessKey string
	SecretKey string
	Addr      string
}

func (d *Docker) MinioConfigs(ctx context.Context) ([]MinioConfig, error) {
	containers, err := d.client.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("status", "running")),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get container list: %w", err)
	}

	configs := make([]MinioConfig, 0, len(containers))
	for _, ctr := range containers {
		if ctr.Image != d.imageName {
			continue
		}
		cfg, err := d.minioConfigFromContainer(ctx, ctr)
		if err != nil {
			return nil, err
		}
		configs = append(configs, cfg)
	}
	return configs, nil
}

func (d *Docker) minioConfigFromContainer(ctx context.Context, ctr types.Container) (MinioConfig, error) {
	networkKey := "distributed-object-storage_amazin-object-storage"
	settings := ctr.NetworkSettings.Networks[networkKey]
	if settings == nil {
		return MinioConfig{}, fmt.Errorf("failed to find %s", networkKey)
	}
	inspect, err := d.client.ContainerInspect(ctx, ctr.ID)
	if err != nil {
		return MinioConfig{}, fmt.Errorf("failed to inspect container: %w", err)
	}
	accessKey, err := valueFromEnvConfig(inspect.Config.Env, "MINIO_ACCESS_KEY")
	if err != nil {
		return MinioConfig{}, err
	}
	secretKey, err := valueFromEnvConfig(inspect.Config.Env, "MINIO_SECRET_KEY")
	if err != nil {
		return MinioConfig{}, err
	}
	cfg := MinioConfig{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Addr:      settings.IPAddress + ":9000",
	}
	return cfg, nil
}

func valueFromEnvConfig(envs []string, name string) (string, error) {
	prefix := name + "="
	idx := slices.IndexFunc(envs, func(s string) bool {
		return strings.HasPrefix(s, prefix)
	})
	if idx == -1 {
		return "", fmt.Errorf("failed to find env in config: %s", name)
	}
	return strings.TrimPrefix(envs[idx], prefix), nil
}
