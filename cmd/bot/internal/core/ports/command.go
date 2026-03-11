package ports

import (
	"context"

	"github.com/thebeyond-net/control-plane/cmd/bot/internal/core/application/input"
	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

type CommandHandler interface {
	Execute(ctx context.Context, msg input.Message, user domain.User) error
}
