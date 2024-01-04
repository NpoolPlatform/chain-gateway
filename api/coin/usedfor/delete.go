//nolint:nolintlint,dupl
package coinusedfor

import (
	"context"

	coinusedfor1 "github.com/NpoolPlatform/chain-gateway/pkg/coin/usedfor"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/usedfor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteCoinUsedFor(ctx context.Context, in *npool.DeleteCoinUsedForRequest) (*npool.DeleteCoinUsedForResponse, error) {
	handler, err := coinusedfor1.NewHandler(
		ctx,
		coinusedfor1.WithID(&in.ID, true),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCoinUsedFor",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCoinUsedForResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.DeleteCoinUsedFor(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"DeleteCoinUsedFor",
			"In", in,
			"Error", err,
		)
		return &npool.DeleteCoinUsedForResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteCoinUsedForResponse{
		Info: info,
	}, nil
}
