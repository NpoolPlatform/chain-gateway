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

func (s *Server) UpdateCoin(ctx context.Context, in *npool.UpdateCoinRequest) (*npool.UpdateCoinResponse, error) {
	handler, err := appcoin1.NewHandler(
		ctx,
		appcoin1.WithID(&in.ID),
		appcoin1.WithName(in.Name),
		appcoin1.WithDisplayNames(in.DisplayNames),
		appcoin1.WithLogo(in.Logo),
		appcoin1.WithForPay(in.ForPay),
		appcoin1.WithWithdrawAutoReviewAmount(in.WithdrawAutoReviewAmount),
		appcoin1.WithMarketValue(in.MarketValue),
		appcoin1.WithSettlePercent(in.SettlePercent),
		appcoin1.WithSettleTips(in.SettleTips),
		appcoin1.WithDailyRewardAmount(in.DailyRewardAmount),
		appcoin1.WithProductPage(in.ProductPage),
		appcoin1.WithDisabled(in.Disabled),
		appcoin1.WithDisplay(in.Display),
		appcoin1.WithDisplayIndex(in.DisplayIndex),
		appcoin1.WithMaxAmountPerWithdraw(in.MaxAmountPerWithdraw),
	)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCoin",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	info, err := handler.UpdateCoin(ctx)
	if err != nil {
		logger.Sugar().Errorw(
			"UpdateCoin",
			"In", in,
			"Error", err,
		)
		return &npool.UpdateCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCoinResponse{
		Info: info,
	}, nil
}
