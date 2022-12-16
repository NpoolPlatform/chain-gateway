package fiatcurrency

import (
	"context"

	"github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/fiatcurrency"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	fiatcurrency.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	fiatcurrency.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return fiatcurrency.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
