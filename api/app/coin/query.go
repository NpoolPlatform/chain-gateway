//nolint:nolintlint,dupl
package appcoin

import (
	"context"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	appcoin1 "github.com/NpoolPlatform/chain-gateway/pkg/appcoin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/appcoin"
	"github.com/google/uuid"
)

func (s *Server) GetCoins(ctx context.Context, in *npool.GetCoinsRequest) (*npool.GetCoinsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetCoins")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))
	span = commontracer.TraceInvoker(span, "coin", "coin", "Rows")

	if _, err := uuid.Parse(in.GetAppID()); err != nil {
		return &npool.GetCoinsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := appcoin1.GetAppCoins(ctx, in.GetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetCoins", "error", err)
		return &npool.GetCoinsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCoinsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCoins(ctx context.Context, in *npool.GetAppCoinsRequest) (*npool.GetAppCoinsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppCoins")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceOffsetLimit(span, int(in.GetOffset()), int(in.GetLimit()))
	span = commontracer.TraceInvoker(span, "coin", "coin", "Rows")

	if _, err := uuid.Parse(in.GetTargetAppID()); err != nil {
		return &npool.GetAppCoinsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	infos, total, err := appcoin1.GetAppCoins(ctx, in.GetTargetAppID(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAppCoins", "error", err)
		return &npool.GetAppCoinsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCoinsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
