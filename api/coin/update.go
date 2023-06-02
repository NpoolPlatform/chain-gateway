//nolint:nolintlint,dupl
package coin

import (
	"context"

	coin1 "github.com/NpoolPlatform/chain-gateway/pkg/coin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateCoin(ctx context.Context, in *npool.UpdateCoinRequest) (*npool.UpdateCoinResponse, error) {
	handler, err := coin1.NewHandler(
		ctx,
		coin1.WithID(&in.ID),
		coin1.WithLogo(in.Logo),
		coin1.WithPresale(in.Presale),
		coin1.WithReservedAmount(in.ReservedAmount),
		coin1.WithForPay(in.ForPay),
		coin1.WithHomePage(in.HomePage),
		coin1.WithSpecs(in.Specs),
		coin1.WithFeeCoinTypeID(in.FeeCoinTypeID),
		coin1.WithWithdrawFeeByStableUSD(in.WithdrawFeeByStableUSD),
		coin1.WithWithdrawFeeAmount(in.WithdrawFeeAmount),
		coin1.WithCollectFeeAmount(in.CollectFeeAmount),
		coin1.WithHotWalletFeeAmount(in.HotWalletFeeAmount),
		coin1.WithHotWalletAccountAmount(in.HotWalletAccountAmount),
		coin1.WithLowFeeAmount(in.LowFeeAmount),
		coin1.WithHotLowFeeAmount(in.HotLowFeeAmount),
		coin1.WithPaymentAccountCollectAmount(in.PaymentAccountCollectAmount),
		coin1.WithDisabled(in.Disabled),
		coin1.WithStableUSD(in.StableUSD),
		coin1.WithLeastTransferAmount(in.LeastTransferAmount),
		coin1.WithNeedMemo(in.NeedMemo),
		coin1.WithRefreshCurrency(in.RefreshCurrency),
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
