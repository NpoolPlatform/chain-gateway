package appcoin

import (
	"context"

	"github.com/NpoolPlatform/message/npool/chain/gw/v1/appcoin"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	appcoin.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	appcoin.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return appcoin.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
