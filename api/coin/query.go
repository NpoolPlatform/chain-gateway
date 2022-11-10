//nolint:nolintlint,dupl
package coin

import (
	"context"

	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/coin"
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

	infos, total, err := coinmwcli.GetCoins(ctx, in.GetConds(), in.GetOffset(), in.GetLimit())
	if err != nil {
		logger.Sugar().Errorf("fail get coins: %v", err)
		return &npool.GetCoinsResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetCoinsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
