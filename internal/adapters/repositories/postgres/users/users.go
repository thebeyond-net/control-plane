package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/thebeyond-net/control-plane/internal/adapters/repositories/postgres/generated"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
	"go.uber.org/zap"
)

type UserRepository struct {
	pool      *pgxpool.Pool
	queries   *sqlc.Queries
	appLogger *zap.Logger
}

func NewRepository(pool *pgxpool.Pool, appLogger *zap.Logger) ports.UserRepository {
	queries := sqlc.New(pool)
	return &UserRepository{pool, queries, appLogger}
}

func (r *UserRepository) GetByIdentity(ctx context.Context, provider, providerID string) (domain.User, error) {
	user, err := r.queries.GetUserByIdentity(ctx, sqlc.GetUserByIdentityParams{
		Provider:   provider,
		ProviderID: providerID,
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user by identity: %w", err)
	}

	var referrerID *string
	if user.ReferrerID.Valid {
		referrer := uuid.UUID(user.ReferrerID.Bytes).String()
		referrerID = &referrer
	}

	return domain.User{
		ID: user.ID.String(),
		Identity: domain.Identity{
			Provider:   provider,
			ProviderID: providerID,
		},
		Subscription: domain.Subscription{
			ExpiresAt: user.SubscriptionExpiresAt.Time,
		},
		Referral: domain.Referral{
			ReferrerID:     referrerID,
			Balance:        int(user.ReferralBalance),
			CommissionRate: int(user.ReferralCommissionRate),
			Count:          int(user.ReferralsCount),
		},
		Devices:      int(user.Devices),
		Bandwidth:    int(user.Bandwidth),
		LanguageCode: user.LanguageCode,
		CurrencyCode: user.CurrencyCode,
		CreatedAt:    user.CreatedAt.Time,
	}, nil
}

func (r *UserRepository) UpdateLanguage(ctx context.Context, userID string, languageCode string) error {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	return r.queries.SetUserLanguageCode(ctx, sqlc.SetUserLanguageCodeParams{
		ID:           uuid,
		LanguageCode: languageCode,
	})
}

func (r *UserRepository) UpdateCurrency(ctx context.Context, userID string, currencyCode string) error {
	uuid, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	return r.queries.SetUserCurrencyCode(ctx, sqlc.SetUserCurrencyCodeParams{
		ID:           uuid,
		CurrencyCode: currencyCode,
	})
}

func (r *UserRepository) CreateWithIdentity(ctx context.Context, user domain.User, identity domain.Identity) error {
	newUserID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txn := r.queries.WithTx(tx)

	referrerID := pgtype.UUID{Valid: false}
	if user.Referral.ReferrerID != nil {
		parsed, err := uuid.Parse(*user.Referral.ReferrerID)
		if err == nil && parsed != uuid.Nil {
			referrerID = pgtype.UUID{Bytes: parsed, Valid: true}
		}
	}

	if err := txn.CreateUser(ctx, sqlc.CreateUserParams{
		ID:        newUserID,
		Devices:   int32(user.Devices),
		Bandwidth: int32(user.Bandwidth),
		SubscriptionExpiresAt: pgtype.Timestamptz{
			Time:  user.Subscription.ExpiresAt,
			Valid: true,
		},
		LanguageCode: user.LanguageCode,
		CurrencyCode: user.CurrencyCode,
		ReferrerID:   referrerID,
	}); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if referrerID.Valid {
		if err := txn.IncrementReferralsCount(ctx, referrerID.Bytes); err != nil {
			return fmt.Errorf("failed to increment referrals count: %w", err)
		}

		if err := txn.AddBonusDays(ctx, referrerID.Bytes); err != nil {
			r.appLogger.Warn("failed to add bonus days to referrer",
				zap.String("referrer_id", referrerID.String()),
				zap.Error(err))
		}
	}

	if err := txn.CreateIdentity(ctx, sqlc.CreateIdentityParams{
		Provider:   identity.Provider,
		ProviderID: identity.ProviderID,
		UserID:     pgtype.UUID{Bytes: newUserID, Valid: true},
	}); err != nil {
		return fmt.Errorf("failed to create identity: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
