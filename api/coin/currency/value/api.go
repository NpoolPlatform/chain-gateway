package currencyvalue

import (
	"context"

	"github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/currency/value"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	value.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	value.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return value.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
