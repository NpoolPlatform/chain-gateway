package chain

import (
	"context"

	chainmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/chain"
	chainmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/chain"
)

func (h *Handler) UpdateChain(ctx context.Context) (*chainmwpb.Chain, error) {
	if err := chainmwcli.UpdateChain(ctx, &chainmwpb.ChainReq{
		ID:         h.ID,
		ChainType:  h.ChainType,
		NativeUnit: h.NativeUnit,
		AtomicUnit: h.AtomicUnit,
		UnitExp:    h.UnitExp,
		ENV:        h.ENV,
		ChainID:    h.ChainID,
		NickName:   h.Nickname,
		GasType:    h.GasType,
		Logo:       h.Logo,
	}); err != nil {
		return nil, err
	}

	return h.GetChain(ctx)
}
