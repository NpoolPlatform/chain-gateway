package api

import (
	"context"

	chaingw "github.com/NpoolPlatform/message/npool/chain/gw/v1"

	"github.com/NpoolPlatform/chain-gateway/api/appcoin"
	"github.com/NpoolPlatform/chain-gateway/api/appcoin/description"
	"github.com/NpoolPlatform/chain-gateway/api/coin"
	feed "github.com/NpoolPlatform/chain-gateway/api/coin/currency/feed"
	value "github.com/NpoolPlatform/chain-gateway/api/coin/currency/value"
	"github.com/NpoolPlatform/chain-gateway/api/tx"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	chaingw.UnimplementedGatewayServer
}

func Register(server grpc.ServiceRegistrar) {
	chaingw.RegisterGatewayServer(server, &Server{})
	coin.Register(server)
	feed.Register(server)
	value.Register(server)
	appcoin.Register(server)
	tx.Register(server)
	description.Register(server)
}

func RegisterGateway(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := chaingw.RegisterGatewayHandlerFromEndpoint(context.Background(), mux, endpoint, opts); err != nil {
		return err
	}
	if err := coin.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := feed.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := value.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := appcoin.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := tx.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	if err := description.RegisterGateway(mux, endpoint, opts); err != nil {
		return err
	}
	return nil
}
