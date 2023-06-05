package tx

import (
	"context"

	tx1 "github.com/NpoolPlatform/chain-gateway/pkg/tx"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/tx"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetTxs(ctx context.Context, in *npool.GetTxsRequest) (*npool.GetTxsResponse, error) {
	handler, err := tx1.NewHandler(
		ctx,
		tx1.WithOffset(in.GetOffset()),
		tx1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetTxs",
			"In", in,
			"Error", err,
		)
		return &npool.GetTxsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetTxs(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetTxs",
			"In", in,
			"Error", err,
		)
		return &npool.GetTxsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetTxsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
