package nodepool

import (
	"context"
	"net/http"
	"sync"

	"github.com/thebeyond-net/control-plane/internal/adapters/vpnclient"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type NodePool struct {
	httpClient *http.Client
	repo       ports.NodeRepository
	mu         sync.RWMutex
	nodes      map[string]ports.VPNClient
	authSecret string
}

func New(repo ports.NodeRepository, authSecret string) ports.NodePool {
	nodes := make(map[string]ports.VPNClient)
	return &NodePool{
		httpClient: http.DefaultClient,
		repo:       repo,
		nodes:      nodes,
		authSecret: authSecret,
	}
}

func (a *NodePool) GetVPNClient(ctx context.Context, nodeID string) (ports.VPNClient, error) {
	a.mu.RLock()
	client, ok := a.nodes[nodeID]
	a.mu.RUnlock()

	if ok {
		return client, nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if client, ok := a.nodes[nodeID]; ok {
		return client, nil
	}

	node, err := a.repo.Get(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	newClient := vpnclient.New(a.httpClient, "http://"+node.Address+":4080", a.authSecret)
	a.nodes[nodeID] = newClient

	return newClient, nil
}
