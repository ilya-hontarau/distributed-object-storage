package gateway

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ilya-hontarau/distributed-object-storage/internal/mock"
)

func TestGateway_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	node1 := mock.NewMockStorageNode(ctrl)
	node1.EXPECT().Upload(context.Background(), "node1", strings.NewReader("node1"), 5)
	node2 := mock.NewMockStorageNode(ctrl)
	node2.EXPECT().Upload(context.Background(), "test1", strings.NewReader("node2"), 5)
	node3 := mock.NewMockStorageNode(ctrl)
	node3.EXPECT().Upload(context.Background(), "test", strings.NewReader("node3"), 5)

	gateway := NewGateway([]StorageNode{node1, node2, node3})
	err := gateway.Upload(context.Background(), "node1", strings.NewReader("node1"), 5)
	assert.NoError(t, err)
	err = gateway.Upload(context.Background(), "test1", strings.NewReader("node2"), 5)
	assert.NoError(t, err)
	err = gateway.Upload(context.Background(), "test", strings.NewReader("node3"), 5)
	assert.NoError(t, err)
}

func TestHashDistribution(t *testing.T) {
	counters := make(map[int]int)
	for range 10000 {
		id := uuid.NewString()
		idx := hash(id, 3)
		counters[idx]++
	}
	assert.Greater(t, counters[0], 3000)
	assert.Greater(t, counters[1], 3000)
	assert.Greater(t, counters[2], 3000)
}
