package currencyfeed

import (
	"context"

	feedmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency/feed"
	feedmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/currency/feed"
)

func (h *Handler) UpdateFeed(ctx context.Context) (*feedmwpb.Feed, error) {
	return feedmwcli.UpdateFeed(ctx, &feedmwpb.FeedReq{
		ID:           h.ID,
		FeedCoinName: h.FeedCoinName,
		Disabled:     h.Disabled,
	})
}
