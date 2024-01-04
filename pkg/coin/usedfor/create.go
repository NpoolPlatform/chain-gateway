package coinusedfor

import (
	"context"

	coinusedformwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/usedfor"
	coinusedformwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/usedfor"
)

func (h *Handler) CreateCoinUsedFor(ctx context.Context) (*coinusedformwpb.CoinUsedFor, error) {
	return coinusedformwcli.CreateCoinUsedFor(ctx, &coinusedformwpb.CoinUsedForReq{
		CoinTypeID: h.CoinTypeID,
		UsedFor:    h.UsedFor,
	})
}
