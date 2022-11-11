//nolint:nolintlint,dupl
package appcoin

import (
	"context"
	"fmt"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/appcoin"
)

func (s *Server) DeleteCoin(ctx context.Context, in *npool.DeleteCoinRequest) (*npool.DeleteCoinResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "DeleteCoin")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if app, err := appmwcli.GetApp(ctx, in.GetTargetAppID()); err != nil || app == nil {
		return &npool.DeleteCoinResponse{}, status.Error(codes.InvalidArgument, fmt.Sprintf("appid is invalid: %v", err))
	}
	ac, err := appcoinmwcli.GetCoin(ctx, in.GetID())
	if err != nil {
		return &npool.DeleteCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if ac.AppID != in.GetTargetAppID() {
		return &npool.DeleteCoinResponse{}, status.Error(codes.InvalidArgument, "permission denied")
	}

	span = commontracer.TraceID(span, in.GetID())
	span = commontracer.TraceInvoker(span, "coin", "coin", "Delete")

	info, err := appcoinmwcli.DeleteCoin(ctx, in.GetID())
	if err != nil {
		return &npool.DeleteCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.DeleteCoinResponse{
		Info: info,
	}, nil
}
