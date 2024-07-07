package chain

import (
	"context"
	"fmt"

	chainmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/chain"
	chainmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/chain"
)

type queryHandler struct {
	*Handler
}

func (h *Handler) GetChains(ctx context.Context) ([]*chainmwpb.Chain, uint32, error) {
	return chainmwcli.GetChains(ctx, &chainmwpb.Conds{}, h.Offset, h.Limit)
}

func (h *Handler) GetChain(ctx context.Context) (*chainmwpb.Chain, error) {
	chain, err := chainmwcli.GetChain(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if chain == nil {
		return nil, fmt.Errorf("invalid chain")
	}

	return chain, nil
}
