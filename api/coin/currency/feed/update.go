//nolint:nolintlint,dupl
package currencyfeed

import (
	"context"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency/feed"
	coinmgrpb "github.com/NpoolPlatform/message/npool/chain/mgr/v1/coin/currency/feed"

	apifeed "github.com/NpoolPlatform/chain-middleware/api/coin/currency/feed"
	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin/currency/feed"
)

func (s *Server) UpdateCurrencyFeed(ctx context.Context, in *npool.UpdateCurrencyFeedRequest) (*npool.UpdateCurrencyFeedResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetCoins")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	err = apifeed.ValidateUpdate(&coinmgrpb.CurrencyFeedReq{
		ID:         &in.ID,
		FeedSource: &in.FeedSource,
	})
	if err != nil {
		return nil, err
	}

	span = commontracer.TraceInvoker(span, "coin", "coin", "CreateCurrencyFeed")
	info, err := coinmwcli.UpdateCurrencyFeed(ctx, &coinmgrpb.CurrencyFeedReq{
		ID:         &in.ID,
		FeedSource: &in.FeedSource,
	})
	if err != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.UpdateCurrencyFeedResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.UpdateCurrencyFeedResponse{
		Info: info,
	}, nil
}
