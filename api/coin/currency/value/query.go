//nolint:nolintlint,dupl
package currencyvalue

import (
	"context"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency/value"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/currency/value"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/currency/value"

	npoolpb "github.com/NpoolPlatform/message/npool"
)

func (s *Server) GetCurrencies(ctx context.Context, in *npool.GetCurrenciesRequest) (*npool.GetCurrenciesResponse, error) {
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

	infos, total, err := coinmwcli.GetCurrencies(ctx, nil, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.GetCurrenciesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCurrenciesResponse{
		Infos: infos,
		Total: total,
	}, nil
}

func (s *Server) GetHistories(ctx context.Context, in *npool.GetHistoriesRequest) (*npool.GetHistoriesResponse, error) {
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

	conds := &coinmwpb.Conds{
		StartAt: &npoolpb.Uint32Val{
			Op:    cruder.EQ,
			Value: in.StartAt,
		},
		EndAt: &npoolpb.Uint32Val{
			Op:    cruder.EQ,
			Value: in.EndAt,
		},
	}

	if in.CoinTypeID != nil {
		conds.CoinTypeID = &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: *in.CoinTypeID,
		}
	}
	infos, total, err := coinmwcli.GetHistories(ctx, conds, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.GetHistoriesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetHistoriesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
