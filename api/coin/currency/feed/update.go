//nolint:nolintlint,dupl
package feed

import (
	"context"

	feed1 "github.com/NpoolPlatform/chain-gateway/pkg/coin/currency/feed"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/currency/feed"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateFeed(ctx context.Context, in *npool.UpdateFeedRequest) (*npool.UpdateFeedResponse, error) {
	handler, err := feed1.NewHandler(
		ctx,
		feed1.WithID(&in.ID),
		feed1.WithFeedCoinName(in.FeedCoinName),
		feed1.WithDisabled(in.Disabled),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateFeed",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateFeedResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateFeed(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateFeed",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateFeedResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateFeedResponse{
		Info: info,
	}, nil
}
