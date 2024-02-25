package gateway

import (
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
)

var ErrNotFound = errors.New("not found")

type StorageNode interface {
	Upload(ctx context.Context, id string, file io.Reader, contentLength int) error
	Download(ctx context.Context, id string) (io.Reader, error)
}

type Gateway struct {
	nodes []StorageNode
}

func NewGateway(nodes []StorageNode) *Gateway {
	return &Gateway{nodes: nodes}
}

func (g *Gateway) Upload(ctx context.Context, id string, file io.Reader, contentLength int) error {
	idx := g.nodeIdx(id)
	node := g.nodes[idx]
	err := node.Upload(ctx, id, file, contentLength)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

func (g *Gateway) Download(ctx context.Context, id string) (io.Reader, error) {
	idx := g.nodeIdx(id)
	node := g.nodes[idx]
	file, err := node.Download(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}
	return file, nil
}

func (g *Gateway) nodeIdx(id string) int {
	return hash(id, len(g.nodes))
}

func hash(id string, maxIdx int) int {
	// TODO: add hash function configurable
	h := fnv.New32a()
	_, _ = h.Write([]byte(id))
	sum32 := h.Sum32()
	u := sum32 % uint32(maxIdx)
	return int(u)
}
