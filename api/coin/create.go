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

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin"
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

	err = coinmw.ValidateCreate(in.GetInfo())
	if err != nil {
		return &npool.CreateCoinResponse{}, err
	}

	span = commontracer.TraceInvoker(span, "coin", "coin", "Create")

	info, err := coinmwcli.CreateCoin(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("fail create coin: %v", err.Error())
		return &npool.CreateCoinResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCoinResponse{
		Info: info,
	}, nil
}
