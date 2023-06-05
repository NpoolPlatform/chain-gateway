package description

import (
	"context"

	descmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin/description"
	descmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin/description"
)

func (h *Handler) CreateCoinDescription(ctx context.Context) (*descmwpb.CoinDescription, error) {
	return descmwcli.CreateCoinDescription(ctx, &descmwpb.CoinDescriptionReq{
		AppID:      h.AppID,
		CoinTypeID: h.CoinTypeID,
		UsedFor:    h.UsedFor,
		Title:      h.Title,
		Message:    h.Message,
	})
}
