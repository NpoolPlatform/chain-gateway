package currencyfeed

import (
	"context"

	feedmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/fiat/currency/feed"
	feedmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/fiat/currency/feed"
)

func (h *Handler) CreateFeed(ctx context.Context) (*feedmwpb.Feed, error) {
	return feedmwcli.CreateFeed(ctx, &feedmwpb.FeedReq{
		FiatID:       h.FiatID,
		FeedType:     h.FeedType,
		FeedFiatName: h.FeedFiatName,
	})
}
