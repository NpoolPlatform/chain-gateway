//nolint:nolintlint,dupl
package appcoin

import (
	"context"

	appcoin1 "github.com/NpoolPlatform/chain-gateway/pkg/app/coin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/app/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateCoin(ctx context.Context, in *npool.CreateCoinRequest) (*npool.CreateCoinResponse, error) {
	handler, err := appcoin1.NewHandler(
		ctx,
		appcoin1.WithAppID(&in.TargetAppID, true),
		appcoin1.WithCoinTypeID(&in.CoinTypeID, true),
		appcoin1.WithName(&in.Name, true),
		appcoin1.WithLogo(&in.Logo, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCoin",
			"In", in,
			"Error", err,
		)
		return &npool.CreateCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.CreateCoin(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"CreateCoin",
			"In", in,
			"Error", err,
		)
		return &npool.CreateCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCoinResponse{
		Info: info,
	}, nil
}
