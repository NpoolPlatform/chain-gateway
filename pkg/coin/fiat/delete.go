package coinfiat

import (
	"context"
	"fmt"

	coinfiatmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/fiat"
	coinfiatmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/fiat"
)

func (h *Handler) DeleteCoinFiat(ctx context.Context) (*coinfiatmwpb.CoinFiat, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid coinfiatid")
	}
	return coinfiatmwcli.DeleteCoinFiat(ctx, *h.ID)
}
