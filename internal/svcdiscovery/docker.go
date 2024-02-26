package svcdiscovery

import (
	"context"
	"fmt"
	"slices"
	"strings"

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
		settings := ctr.NetworkSettings.Networks["distributed-object-storage_amazin-object-storage"]
		if settings == nil {
			// TODO
		}
		inspect, err := d.client.ContainerInspect(ctx, ctr.ID)
		if err != nil {
			return nil, err
		}
		accessKeyIdx := slices.IndexFunc(inspect.Config.Env, func(s string) bool {
			return strings.HasPrefix(s, "MINIO_ACCESS_KEY=")
		})
		if accessKeyIdx == -1 {
			// TODO
		}
		secretKeyIdx := slices.IndexFunc(inspect.Config.Env, func(s string) bool {
			return strings.HasPrefix(s, "MINIO_SECRET_KEY=")
		})
		if secretKeyIdx == -1 {
			// TODO
		}
		configs = append(configs, MinioConfig{
			AccessKey: strings.TrimPrefix(inspect.Config.Env[accessKeyIdx], "MINIO_ACCESS_KEY="),
			SecretKey: strings.TrimPrefix(inspect.Config.Env[secretKeyIdx], "MINIO_SECRET_KEY="),
			Addr:      settings.IPAddress + ":9000",
		})
	}
	return configs, nil
}
