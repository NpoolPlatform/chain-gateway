package fiat

import (
	"context"

	currency "github.com/NpoolPlatform/message/npool/chain/gw/v1/fiat"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	currency.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	currency.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return currency.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
