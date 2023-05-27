package description

import (
	"context"
	"fmt"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	descmw "github.com/NpoolPlatform/chain-middleware/api/appcoin/description"
	descmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin/description"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/appcoin/description"
	descmgrpb "github.com/NpoolPlatform/message/npool/chain/mgr/v1/appcoin/description"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//nolint
func (s *Server) CreateCoinDescription(
	ctx context.Context,
	in *npool.CreateCoinDescriptionRequest,
) (
	*npool.CreateCoinDescriptionResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateCoinDescription")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	req := &descmgrpb.CoinDescriptionReq{
		AppID:      &in.AppID,
		CoinTypeID: &in.CoinTypeID,
		UsedFor:    &in.UsedFor,
		Title:      &in.Title,
		Message:    &in.Message,
	}

	if err := descmw.ValidateCreate(ctx, req); err != nil {
		logger.Sugar().Errorw("CreateCoinDescription", "error", err)
		return &npool.CreateCoinDescriptionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app, err := appmwcli.GetApp(ctx, in.GetAppID()); err != nil || app == nil {
		logger.Sugar().Errorw("CreateCoinDescription", "AppID", in.GetAppID(), "error", err)
		return &npool.CreateCoinDescriptionResponse{}, status.Error(codes.InvalidArgument, fmt.Sprintf("appid is invalid: %v", err))
	}

	span = commontracer.TraceInvoker(span, "description", "description", "Create")

	info, err := descmwcli.CreateCoinDescription(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("CreateCoinDescription", "error", err)
		return &npool.CreateCoinDescriptionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateCoinDescriptionResponse{
		Info: info,
	}, nil
}

//nolint
func (s *Server) CreateAppCoinDescription(
	ctx context.Context,
	in *npool.CreateAppCoinDescriptionRequest,
) (
	*npool.CreateAppCoinDescriptionResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "CreateAppCoinDescription")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	req := &descmgrpb.CoinDescriptionReq{
		AppID:      &in.TargetAppID,
		CoinTypeID: &in.CoinTypeID,
		UsedFor:    &in.UsedFor,
		Title:      &in.Title,
		Message:    &in.Message,
	}

	if err := descmw.ValidateCreate(ctx, req); err != nil {
		logger.Sugar().Errorw("CreateAppCoinDescription", "error", err)
		return &npool.CreateAppCoinDescriptionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app, err := appmwcli.GetApp(ctx, in.GetTargetAppID()); err != nil || app == nil {
		logger.Sugar().Errorw("CreateAppCoinDescription", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.CreateAppCoinDescriptionResponse{}, status.Error(codes.InvalidArgument, fmt.Sprintf("appid is invalid: %v", err))
	}

	span = commontracer.TraceInvoker(span, "description", "description", "Create")

	info, err := descmwcli.CreateCoinDescription(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("CreateAppCoinDescription", "error", err)
		return &npool.CreateAppCoinDescriptionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateAppCoinDescriptionResponse{
		Info: info,
	}, nil
}
