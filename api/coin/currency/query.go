//nolint:nolintlint,dupl
package currency

import (
	"context"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"

	constant1 "github.com/NpoolPlatform/chain-gateway/pkg/const"
	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"
	coininfocli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	currencymwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency"
	npoolpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/currency"
	currencymwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/currency"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	limit := constant1.DefaultRowLimit
	if in.GetLimit() > 0 {
		limit = in.GetLimit()
	}

	coinInfos, total, err := coininfocli.GetCoins(ctx, nil, in.GetOffset(), limit)
	if err != nil {
		return nil, err
	}
	if len(coinInfos) == 0 {
		return &npool.GetCurrenciesResponse{
			Infos: nil,
			Total: total,
		}, nil
	}

	coinTypeIDs := []string{}
	for _, val := range coinInfos {
		coinTypeIDs = append(coinTypeIDs, val.ID)
	}

	infos, err := currencymwcli.GetCurrencies(ctx, &currencymwpb.Conds{
		CoinTypeIDs: &npoolpb.StringSliceVal{
			Op:    cruder.IN,
			Value: coinTypeIDs,
		},
	})
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

	conds := &currencymwpb.Conds{}

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

	if in.CoinTypeID != nil {
		conds.CoinTypeID = &npoolpb.StringVal{
			Op:    cruder.EQ,
			Value: in.GetCoinTypeID(),
		}
	}
	infos, total, err := currencymwcli.GetHistories(ctx, conds, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.GetHistoriesResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetHistoriesResponse{
		Infos: infos,
		Total: total,
	}, nil
}
