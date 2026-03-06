package userctx

import (
	"context"
	"errors"

	"github.com/thebeyond-net/control-plane/internal/core/domain"
)

type ctxKey struct{}

var ErrNoUserInContext = errors.New("no user in context")

func WithUser(ctx context.Context, user domain.User) context.Context {
	return context.WithValue(ctx, ctxKey{}, user)
}

func FromContext(ctx context.Context) (domain.User, error) {
	value := ctx.Value(ctxKey{})
	if value == nil {
		return domain.User{}, ErrNoUserInContext
	}

	user, ok := value.(domain.User)
	if !ok {
		return domain.User{}, ErrNoUserInContext
	}

	return user, nil
}
