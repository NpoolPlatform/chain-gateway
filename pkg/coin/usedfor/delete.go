package coinusedfor

import (
	"context"
	"fmt"

	coinusedformwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/usedfor"
	coinusedformwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/usedfor"
)

func (h *Handler) DeleteCoinUsedFor(ctx context.Context) (*coinusedformwpb.CoinUsedFor, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid coinusedforid")
	}
	return coinusedformwcli.DeleteCoinUsedFor(ctx, *h.ID)
}
