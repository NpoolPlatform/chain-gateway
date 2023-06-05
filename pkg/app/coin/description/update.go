package description

import (
	"context"

	descmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin/description"
	descmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin/description"
)

func (h *Handler) UpdateCoinDescription(ctx context.Context) (*descmwpb.CoinDescription, error) {
	return descmwcli.UpdateCoinDescription(ctx, &descmwpb.CoinDescriptionReq{
		ID:      h.ID,
		Title:   h.Title,
		Message: h.Message,
	})
}
