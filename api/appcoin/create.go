//nolint:nolintlint,dupl
package appcoin

import (
	"context"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appcoinmw "github.com/NpoolPlatform/chain-middleware/api/appcoin"
	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/appcoin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/appcoin"
)

func (s *Server) CreateCoin(ctx context.Context, in *npool.CreateCoinRequest) (*npool.CreateCoinResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateCoin")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	_, err = appcoinmw.ValidateCreate(ctx, &appcoinmwpb.CoinReq{
		AppID:      &in.TargetAppID,
		CoinTypeID: &in.CoinTypeID,
	})
	if err != nil {
		return &npool.CreateCoinResponse{}, err
	}

	if _, err := appmwcli.GetApp(ctx, in.GetTargetAppID()); err != nil {
		return &npool.CreateCoinResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	span = commontracer.TraceInvoker(span, "coin", "coin", "Create")

	info, err := appcoinmwcli.CreateCoin(ctx, &appcoinmwpb.CoinReq{
		AppID:      &in.TargetAppID,
		CoinTypeID: &in.CoinTypeID,
	})
	if err != nil {
		return &npool.CreateCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCoinResponse{
		Info: info,
	}, nil
}
