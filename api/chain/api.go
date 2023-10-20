package chain

import (
	"context"

	"github.com/NpoolPlatform/message/npool/chain/gw/v1/chain"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	chain.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	chain.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return chain.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
