package fiat

import (
	"context"

	fiatmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/fiat"
	fiatmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/fiat"
)

func (h *Handler) UpdateFiat(ctx context.Context) (*fiatmwpb.Fiat, error) {
	return fiatmwcli.UpdateFiat(ctx, &fiatmwpb.FiatReq{
		ID:   h.ID,
		Name: h.Name,
		Unit: h.Unit,
		Logo: h.Logo,
	})
}
