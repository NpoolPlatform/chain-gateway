package coinusedfor

import (
	"context"
	"fmt"

	coinusedformwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/usedfor"
	coinusedformwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/usedfor"
)

func (h *Handler) DeleteCoinUsedFor(ctx context.Context) (*coinusedformwpb.CoinUsedFor, error) {
	if h.ID == nil || h.EntID == nil {
		return nil, fmt.Errorf("invalid coinusedforid")
	}
	info, err := coinusedformwcli.GetCoinUsedFor(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("invalid coinusedfor")
	}
	if info.ID != *h.ID {
		return nil, fmt.Errorf("invalid coinusedforid")
	}
	return coinusedformwcli.DeleteCoinUsedFor(ctx, *h.ID)
}
