package vpnclient

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/thebeyond-net/control-plane/internal/core/ports"
	v1 "github.com/thebeyond-net/node-agent/pkg/amneziawg/v1"
	"github.com/thebeyond-net/node-agent/pkg/amneziawg/v1/amneziawgv1connect"
)

type Adapter struct {
	client amneziawgv1connect.AmneziaWGServiceClient
}

func New(httpClient *http.Client, baseURL, authSecret string) ports.VPNClient {
	interceptor := connect.WithInterceptors(
		authInterceptor(authSecret),
	)

	return &Adapter{
		client: amneziawgv1connect.NewAmneziaWGServiceClient(
			httpClient,
			baseURL,
			interceptor,
		),
	}
}

func (a *Adapter) CreatePeer(ctx context.Context, bandwidth int) (string, string, error) {
	res, err := a.client.CreatePeer(ctx, &v1.CreatePeerRequest{
		Bandwidth: int32(bandwidth),
	})
	if err != nil {
		return "", "", err
	}
	return res.Pubkey, res.Config, nil
}

func (a *Adapter) DeletePeer(ctx context.Context, pubkey string) error {
	_, err := a.client.DeletePeer(ctx, &v1.DeletePeerRequest{
		Pubkey: pubkey,
	})
	return err
}

func authInterceptor(secret string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			req.Header().Set("Authorization", secret)
			return next(ctx, req)
		}
	}
}
