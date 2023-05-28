package currencyfeed

import (
	"context"

	feedmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency/feed"
	feedmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/currency/feed"
)

func (h *Handler) CreateFeed(ctx context.Context) (*feedmwpb.Feed, error) {
	return feedmwcli.CreateFeed(ctx, &feedmwpb.FeedReq{
		CoinTypeID:   h.CoinTypeID,
		FeedType:     h.FeedType,
		FeedCoinName: h.FeedCoinName,
	})
}
