package fiat

import (
	"context"

	fiatmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/fiat"
	fiatmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/fiat"
)

func (h *Handler) GetFiats(ctx context.Context) ([]*fiatmwpb.Fiat, uint32, error) {
	return fiatmwcli.GetFiats(ctx, &fiatmwpb.Conds{}, h.Offset, h.Limit)
}
