package nodes

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/thebeyond-net/control-plane/internal/adapters/repositories/postgres/generated"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type NodeRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) ports.NodeRepository {
	queries := sqlc.New(pool)
	return &NodeRepository{pool, queries}
}

func (r *NodeRepository) Get(ctx context.Context, id string) (domain.Node, error) {
	node, err := r.queries.GetNode(ctx, id)
	return domain.Node{
		ID:      node.ID,
		Address: node.Address,
	}, err
}

func (r *NodeRepository) List(ctx context.Context) ([]domain.Node, error) {
	rows, err := r.queries.ListNodes(ctx)
	if err != nil {
		return nil, err
	}

	nodes := make([]domain.Node, len(rows))
	for i, node := range rows {
		nodes[i] = domain.Node{
			ID:      node.ID,
			Address: node.Address,
		}
	}

	return nodes, nil
}
