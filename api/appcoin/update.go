//nolint:nolintlint,dupl
package appcoin

import (
	"context"
	"fmt"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	usermwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/user"

	appcoinmw "github.com/NpoolPlatform/chain-middleware/api/appcoin"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/appcoin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/appcoin"
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

	req := &appcoinmwpb.CoinReq{
		ID:                       &in.ID,
		AppID:                    &in.AppID,
		CoinTypeID:               &in.CoinTypeID,
		Name:                     in.Name,
		Logo:                     in.Logo,
		ForPay:                   in.ForPay,
		WithdrawAutoReviewAmount: in.WithdrawAutoReviewAmount,
		DailyRewardAmount:        in.DailyRewardAmount,
		MarketValue:              in.MarketValue,
		SettlePercent:            in.SettlePercent,
		Setter:                   &in.UserID,
		Disabled:                 in.Disabled,
		ProductPage:              in.ProductPage,
		Display:                  in.Display,
	}

	if err := appcoinmw.ValidateUpdate(ctx, req); err != nil {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if app, err := appmwcli.GetApp(ctx, in.GetAppID()); err != nil || app == nil {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, fmt.Sprintf("appid is invalid: %v", err))
	}
	ac, err := appcoinmwcli.GetCoin(ctx, in.GetID())
	if err != nil {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if ac.AppID != in.GetAppID() {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, "permission denied")
	}
	if ac.CoinDisabled && !in.GetDisabled() {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, "can not set Disabled to false when CoinDisabled is true")
	}
	if !ac.CoinForPay && in.GetForPay() {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, "can not set ForPay to true when CoinForPay is false")
	}
	user, err := usermwcli.GetUser(ctx, in.GetAppID(), in.GetUserID())
	if err != nil {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if user == nil {
		return &npool.UpdateCoinResponse{}, status.Error(codes.InvalidArgument, "UserID is invalid")
	}

	span = commontracer.TraceID(span, in.GetID())
	span = commontracer.TraceInvoker(span, "coin", "coin", "Update")

	info, err := appcoinmwcli.UpdateCoin(ctx, req)
	if err != nil {
		return &npool.UpdateCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCoinResponse{
		Info: info,
	}, nil
}
