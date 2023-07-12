//nolint:dupl
package description

import (
	"context"

	description1 "github.com/NpoolPlatform/chain-gateway/pkg/app/coin/description"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/app/coin/description"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCoinDescriptions(ctx context.Context, in *npool.GetCoinDescriptionsRequest) (*npool.GetCoinDescriptionsResponse, error) {
	handler, err := description1.NewHandler(
		ctx,
		description1.WithAppID(&in.AppID),
		description1.WithOffset(in.GetOffset()),
		description1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCoinDescriptions",
			"In", in,
			"Error", err,
		)
		return &npool.GetCoinDescriptionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoinDescriptions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetCoinDescriptions",
			"In", in,
			"Error", err,
		)
		return &npool.GetCoinDescriptionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCoinDescriptionsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCoinDescriptions(ctx context.Context, in *npool.GetAppCoinDescriptionsRequest) (*npool.GetAppCoinDescriptionsResponse, error) {
	handler, err := description1.NewHandler(
		ctx,
		description1.WithAppID(&in.TargetAppID),
		description1.WithOffset(in.GetOffset()),
		description1.WithLimit(in.GetLimit()),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppCoinDescriptions",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppCoinDescriptionsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := handler.GetCoinDescriptions(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"GetAppCoinDescriptions",
			"In", in,
			"Error", err,
		)
		return &npool.GetAppCoinDescriptionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCoinDescriptionsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
