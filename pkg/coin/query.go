package coin

import (
	"context"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"
)

func (h *Handler) GetCoins(ctx context.Context) ([]*coinmwpb.Coin, uint32, error) {
	return coinmwcli.GetCoins(ctx, &coinmwpb.Conds{}, h.Offset, h.Limit)
}
