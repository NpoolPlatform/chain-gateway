//nolint:dupl
package description

import (
	"context"
	"fmt"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	descmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin/description"
	commonpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/appcoin/description"
	descmgrpb "github.com/NpoolPlatform/message/npool/chain/mgr/v1/appcoin/description"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetCoinDescriptions(
	ctx context.Context,
	in *npool.GetCoinDescriptionsRequest,
) (
	*npool.GetCoinDescriptionsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetCoinDescriptions")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if app, err := appmwcli.GetApp(ctx, in.GetAppID()); err != nil || app == nil {
		logger.Sugar().Errorw("GetCoinDescriptions", "AppID", in.GetAppID(), "error", err)
		return &npool.GetCoinDescriptionsResponse{}, status.Error(codes.InvalidArgument, fmt.Sprintf("appid is invalid: %v", err))
	}

	span = commontracer.TraceInvoker(span, "description", "description", "Create")

	infos, total, err := descmwcli.GetCoinDescriptions(ctx, &descmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetCoinDescriptions", "error", err)
		return &npool.GetCoinDescriptionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCoinDescriptionsResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetAppCoinDescriptions(
	ctx context.Context,
	in *npool.GetAppCoinDescriptionsRequest,
) (
	*npool.GetAppCoinDescriptionsResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetAppCoinDescriptions")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	if app, err := appmwcli.GetApp(ctx, in.GetTargetAppID()); err != nil || app == nil {
		logger.Sugar().Errorw("GetAppCoinDescriptions", "TargetAppID", in.GetTargetAppID(), "error", err)
		return &npool.GetAppCoinDescriptionsResponse{}, status.Error(codes.InvalidArgument, fmt.Sprintf("appid is invalid: %v", err))
	}

	span = commontracer.TraceInvoker(span, "description", "description", "Create")

	infos, total, err := descmwcli.GetCoinDescriptions(ctx, &descmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetTargetAppID(),
		},
	}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorw("GetAppCoinDescriptions", "error", err)
		return &npool.GetAppCoinDescriptionsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppCoinDescriptionsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
