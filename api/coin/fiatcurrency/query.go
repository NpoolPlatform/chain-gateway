//nolint:nolintlint,dupl
package fiatcurrency

import (
	"context"

	"github.com/google/uuid"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"
	fiatcurrencymwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/fiatcurrency"
	npoolpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/fiatcurrency"
	fiatcurrencymwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/fiatcurrency"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fiatcurrencytypemgrcli "github.com/NpoolPlatform/chain-manager/pkg/client/coin/fiatcurrencytype"
	fiatcurrencytypepb "github.com/NpoolPlatform/message/npool/chain/mgr/v1/coin/fiatcurrencytype"
)

func (s *Server) GetCoinFiatCurrencies(
	ctx context.Context,
	in *npool.GetCoinFiatCurrenciesRequest,
) (
	*npool.GetCoinFiatCurrenciesResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetCoins")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "coin", "coin", "Rows")

	infos, err := fiatcurrencymwcli.GetCoinFiatCurrencies(ctx, in.GetCoinTypeIDs(), in.GetFiatCurrencyTypeIDs())
	if err != nil {
		logger.Sugar().Errorf("fail get ciat currencies: %v", err)
		return &npool.GetCoinFiatCurrenciesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCoinFiatCurrenciesResponse{
		Infos: infos,
	}, nil
}

func (s *Server) GetHistories(
	ctx context.Context,
	in *npool.GetHistoriesRequest,
) (
	*npool.GetHistoriesResponse,
	error,
) {
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

	conds := &fiatcurrencymwpb.Conds{}

	if in.StartAt != nil {
		conds.StartAt = &npoolpb.Uint32Val{
			Op:    cruder.GTE,
			Value: in.GetStartAt(),
		}
	}
	if in.EndAt != nil {
		conds.EndAt = &npoolpb.Uint32Val{
			Op:    cruder.LTE,
			Value: in.GetEndAt(),
		}
	}

	if in.FiatCurrencyTypeID != nil {
		conds.FiatCurrencyTypeID = &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetFiatCurrencyTypeID(),
		}
	}
	infos, total, err := fiatcurrencymwcli.GetHistories(ctx, conds, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.GetHistoriesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetHistoriesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetFiatCurrencyTypes(
	ctx context.Context,
	in *npool.GetFiatCurrencyTypesRequest,
) (*npool.GetFiatCurrencyTypesResponse, error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetCoinFiatCurrencies")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	infos, total, err := fiatcurrencytypemgrcli.GetFiatCurrencyTypes(ctx, nil, in.GetOffset(), in.GetLimit())
	if err != nil {
		return nil, err
	}
	return &npool.GetFiatCurrencyTypesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) CreateFiatCurrencyType(
	ctx context.Context,
	in *npool.CreateFiatCurrencyTypeRequest,
) (
	*npool.CreateFiatCurrencyTypeResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetCoinFiatCurrencies")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()
	if in.GetName() == "" {
		logger.Sugar().Errorf("name is empty")
		return &npool.CreateFiatCurrencyTypeResponse{}, status.Error(codes.InvalidArgument, "name is empty")
	}
	exist, err := fiatcurrencytypemgrcli.ExistFiatCurrencyTypeConds(ctx, &fiatcurrencytypepb.Conds{
		Name: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetName(),
		},
	})
	if err != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.CreateFiatCurrencyTypeResponse{}, status.Error(codes.Internal, err.Error())
	}
	if exist {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.CreateFiatCurrencyTypeResponse{}, status.Error(codes.AlreadyExists, "name already exists")
	}

	info, err := fiatcurrencytypemgrcli.CreateFiatCurrencyType(ctx, &fiatcurrencytypepb.FiatCurrencyTypeReq{
		Name: &in.Name,
	})
	if err != nil {
		return nil, err
	}
	return &npool.CreateFiatCurrencyTypeResponse{
		Info: info,
	}, nil
}

func (s *Server) UpdateFiatCurrencyType(
	ctx context.Context,
	in *npool.UpdateFiatCurrencyTypeRequest,
) (
	*npool.UpdateFiatCurrencyTypeResponse,
	error,
) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetCoinFiatCurrencies")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()
	if in.GetName() == "" {
		logger.Sugar().Errorf("name is empty")
		return &npool.UpdateFiatCurrencyTypeResponse{}, status.Error(codes.InvalidArgument, "name is empty")
	}
	if _, err := uuid.Parse(in.GetID()); err != nil {
		logger.Sugar().Errorf("id is invalid:%v", err)
		return &npool.UpdateFiatCurrencyTypeResponse{}, status.Error(codes.InvalidArgument, "id is invalid")
	}

	fiatCurrencyType, err := fiatcurrencytypemgrcli.GetFiatCurrencyTypeOnly(ctx, &fiatcurrencytypepb.Conds{
		Name: &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetName(),
		},
	})
	if err != nil {
		return nil, err
	}

	if fiatCurrencyType.GetID() != in.GetID() && fiatCurrencyType != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.UpdateFiatCurrencyTypeResponse{}, status.Error(codes.AlreadyExists, "name already exists")
	}

	info, err := fiatcurrencytypemgrcli.UpdateFiatCurrencyType(ctx, &fiatcurrencytypepb.FiatCurrencyTypeReq{
		ID:   &in.ID,
		Name: &in.Name,
	})
	if err != nil {
		return nil, err
	}
	return &npool.UpdateFiatCurrencyTypeResponse{
		Info: info,
	}, nil
}
