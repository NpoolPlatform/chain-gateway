package chain

import (
	"context"

	chainmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/chain"
	chainmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/chain"
)

func (h *Handler) GetChains(ctx context.Context) ([]*chainmwpb.Chain, uint32, error) {
	return chainmwcli.GetChains(ctx, &chainmwpb.Conds{}, h.Offset, h.Limit)
}
