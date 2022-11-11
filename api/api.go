package api

import (
	"context"

	chaingw "github.com/NpoolPlatform/message/npool/chain/gw/v1"

	"github.com/NpoolPlatform/chain-gateway/api/appcoin"
	"github.com/NpoolPlatform/chain-gateway/api/coin"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	chaingw.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	chaingw.RegisterGatewayServer(server, &Server{})
	coin.Register(server)
	appcoin.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := chaingw.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := coin.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appcoin.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
