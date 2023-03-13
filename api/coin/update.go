//nolint:nolintlint,dupl
package coin

import (
	"context"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	coinmw "github.com/NpoolPlatform/chain-middleware/api/coin"
	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"
)

func (s *Server) UpdateCoin(ctx context.Context, in *npool.UpdateCoinRequest) (*npool.UpdateCoinResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateCoin")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	req := &coinmwpb.CoinReq{
		ID:                          &in.ID,
		Presale:                     in.Presale,
		ReservedAmount:              in.ReservedAmount,
		ForPay:                      in.ForPay,
		HomePage:                    in.HomePage,
		Specs:                       in.Specs,
		FeeCoinTypeID:               in.FeeCoinTypeID,
		WithdrawFeeByStableUSD:      in.WithdrawFeeByStableUSD,
		WithdrawFeeAmount:           in.WithdrawFeeAmount,
		CollectFeeAmount:            in.CollectFeeAmount,
		HotWalletFeeAmount:          in.HotWalletFeeAmount,
		LowFeeAmount:                in.LowFeeAmount,
		HotLowFeeAmount:             in.HotLowFeeAmount,
		HotWalletAccountAmount:      in.HotWalletAccountAmount,
		PaymentAccountCollectAmount: in.PaymentAccountCollectAmount,
		Disabled:                    in.Disabled,
		StableUSD:                   in.StableUSD,
		LeastTransferAmount:         in.LeastTransferAmount,
	}

	if err := coinmw.ValidateUpdate(ctx, req); err != nil {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceID(span, in.GetID())
	span = commontracer.TraceInvoker(span, "coin", "coin", "Update")

	info, err := coinmwcli.UpdateCoin(ctx, req)
	if err != nil {
		return &npool.UpdateCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCoinResponse{
		Info: info,
	}, nil
}
