package subscription

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/thebeyond-net/control-plane/internal/adapters/repositories/postgres/generated"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type SubscriptionRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) ports.SubscriptionRepository {
	queries := sqlc.New(pool)
	return &SubscriptionRepository{pool, queries}
}

func (r *SubscriptionRepository) Update(ctx context.Context, userID string, devices, bandwidth, days int) error {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	return r.queries.UpdateSubscription(ctx, sqlc.UpdateSubscriptionParams{
		ID:        uuid,
		Devices:   int32(devices),
		Bandwidth: int32(bandwidth),
		Column3:   int32(days),
	})
}

func (r *SubscriptionRepository) Deactivate(ctx context.Context, userID string) error {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	return r.queries.DeactivateSubscription(ctx, uuid)
}
