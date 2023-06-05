package currencyfeed

import (
	"context"

	feedmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/fiat/currency/feed"
	feedmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/fiat/currency/feed"
)

func (h *Handler) GetFeeds(ctx context.Context) ([]*feedmwpb.Feed, uint32, error) {
	return feedmwcli.GetFeeds(ctx, &feedmwpb.Conds{}, h.Offset, h.Limit)
}
