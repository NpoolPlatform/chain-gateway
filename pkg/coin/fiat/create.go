package coinfiat

import (
	"context"

	coinfiatmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/fiat"
	coinfiatmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/fiat"
)

func (h *Handler) CreateCoinFiat(ctx context.Context) (*coinfiatmwpb.CoinFiat, error) {
	return coinfiatmwcli.CreateCoinFiat(ctx, &coinfiatmwpb.CoinFiatReq{
		CoinTypeID: h.CoinTypeID,
		FiatID:     h.FiatID,
		FeedType:   h.FeedType,
	})
}
