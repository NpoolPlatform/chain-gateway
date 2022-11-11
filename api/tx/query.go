package tx

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/tx"

	constant "github.com/NpoolPlatform/chain-gateway/pkg/message/const"
	commontracer "github.com/NpoolPlatform/chain-gateway/pkg/tracer"

	"go.opentelemetry.io/otel"
	scodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tx1 "github.com/NpoolPlatform/chain-gateway/pkg/tx"
)

func (s *Server) GetTxs(ctx context.Context, in *npool.GetTxsRequest) (*npool.GetTxsResponse, error) {
	var err error

	_, span := otel.Tracer(constant.ServiceName).Start(ctx, "GetTxs")
	defer span.End()

	defer func() {
		if err != nil {
			span.SetStatus(scodes.Error, err.Error())
			span.RecordError(err)
		}
	}()

	span = commontracer.TraceInvoker(span, "tx", "tx", "Create")

	infos, total, err := tx1.GetTxs(ctx, in.GetOffset(), in.GetLimit())
	if err != nil {
		return &npool.GetTxsResponse{}, status.Error(codes.InvalidArgument, err.Error())
	}

	return &npool.GetTxsResponse{
		Infos: infos,
		Total: total,
	}, nil
}
