//nolint:nolintlint,dupl
package currencyfeed

import (
	"context"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency/feed"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/currency/feed"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/currency/feed"
)

func (s *Server) GetCurrencyFeeds(ctx context.Context, in *npool.GetCurrencyFeedsRequest) (*npool.GetCurrencyFeedsResponse, error) {
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

	infos, total, err := coinmwcli.GetCurrencyFeeds(ctx, &coinmwpb.Conds{}, in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.GetCurrencyFeedsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCurrencyFeedsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
