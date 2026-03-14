package devices

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/thebeyond-net/control-plane/internal/adapters/repositories/postgres/generated"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
)

type DeviceRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) ports.DeviceRepository {
	queries := sqlc.New(pool)
	return &DeviceRepository{pool, queries}
}

func (r *DeviceRepository) GetByPublicKey(ctx context.Context, userID, pubkey string) (device domain.Device, err error) {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return domain.Device{}, fmt.Errorf("invalid uuid: %w", err)
	}

	row, err := r.queries.GetDeviceByPublicKey(ctx, sqlc.GetDeviceByPublicKeyParams{
		UserID: uuid,
		Pubkey: pubkey,
	})
	if err != nil {
		return device, err
	}

	return domain.Device{
		PublicKey: row.Pubkey,
		NodeID:    row.NodeID.String,
		Name:      row.Name,
	}, nil
}

func (r *DeviceRepository) List(ctx context.Context, userID string) (devices []domain.Device, err error) {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	rows, err := r.queries.ListDevices(ctx, uuid)
	if err != nil {
		return devices, err
	}

	for _, row := range rows {
		devices = append(devices, domain.Device{
			PublicKey: row.Pubkey,
			NodeID:    row.NodeID.String,
			Name:      row.Name,
		})
	}

	return devices, nil
}

func (r *DeviceRepository) Save(ctx context.Context, device domain.Device) error {
	uuid, err := uuid.Parse(device.UserID)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	return r.queries.CreateDevice(ctx, sqlc.CreateDeviceParams{
		UserID: uuid,
		Pubkey: device.PublicKey,
		NodeID: pgtype.Text{String: device.NodeID, Valid: true},
		Name:   device.Name,
	})
}

func (r *DeviceRepository) Delete(ctx context.Context, userID, pubkey string) error {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	return r.queries.DeleteDevice(ctx, sqlc.DeleteDeviceParams{
		UserID: uuid,
		Pubkey: pubkey,
	})
}
