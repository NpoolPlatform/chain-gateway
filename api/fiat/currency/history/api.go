package currencyhistory

import (
	"context"

	"github.com/NpoolPlatform/message/npool/chain/gw/v1/fiat/currency/history"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	history.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	history.RegisterGatewayServer(server, &Server{})
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return history.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
}
