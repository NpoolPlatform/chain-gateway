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

func (s *Server) GetFeeds(ctx context.Context, in *npool.GetFeedsRequest) (*npool.GetFeedsResponse, error) {
	handler, err := feed1.NewHandler(
		ctx,
		feed1.WithOffset(in.GetOffset()),
		feed1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetFeeds",
			"In", in,
			"Error", err,
		)
		return &npool.GetFeedsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetFeeds(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetFeeds",
			"In", in,
			"Error", err,
		)
		return &npool.GetFeedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetFeedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
