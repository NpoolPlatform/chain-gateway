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

func (s *Server) DeleteCoin(ctx context.Context, in *npool.DeleteCoinRequest) (*npool.DeleteCoinResponse, error) {
	handler, err := appcoin1.NewHandler(
		ctx,
		appcoin1.WithID(&in.ID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCoin",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteCoin(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCoin",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteCoinResponse{
		Info: info,
	}, nil
}
