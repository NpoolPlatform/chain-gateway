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

func (s *Server) UpdateCoinDescription(
	ctx context.Context,
	in *npool.UpdateCoinDescriptionRequest,
) (
	*npool.UpdateCoinDescriptionResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "UpdateCoinDescription")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	req := &descmgrpb.CoinDescriptionReq{
		ID:      &in.ID,
		AppID:   &in.AppID,
		Title:   in.Title,
		Message: in.Message,
	}

	if err := descmw.ValidateUpdate(req); err != nil {
		logger.Sugar().Errorw("UpdateCoinDescription", "error", err)
		return &npool.UpdateCoinDescriptionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if app, err := appmwcli.GetApp(ctx, in.GetAppID()); err != nil || app == nil {
		logger.Sugar().Errorw("UpdateCoinDescription", "AppID", in.GetAppID(), "error", err)
		return &npool.UpdateCoinDescriptionResponse{}, status.Error(codes.InvalidArgument, fmt.Sprintf("appid is invalid: %v", err))
	}
	info, err := descmwcli.GetCoinDescription(ctx, in.GetID())
	if err != nil {
		logger.Sugar().Errorw("UpdateCoinDescription", "error", err)
		return &npool.UpdateCoinDescriptionResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}
	if info.AppID != in.GetAppID() {
		logger.Sugar().Errorw("UpdateCoinDescription", "error", "permission denied")
		return &npool.UpdateCoinDescriptionResponse{}, status.Error(codes.InvalidArgument, "permission denied")
	}

	span = commontracer.TraceInvoker(span, "description", "description", "Update")

	info, err = descmwcli.UpdateCoinDescription(ctx, req)
	if err != nil {
		logger.Sugar().Errorw("UpdateCoinDescription", "error", err)
		return &npool.UpdateCoinDescriptionResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCoinDescriptionResponse{
		Info: info,
	}, nil
}
