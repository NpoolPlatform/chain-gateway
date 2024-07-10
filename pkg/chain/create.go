package chain

import (
	"context"

	chainmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/chain"
	chainmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/chain"

	"github.com/google/uuid"
)

func (h *Handler) CreateChain(ctx context.Context) (*chainmwpb.Chain, error) {
	if h.EntID == nil {
		h.EntID = func() *string { s := uuid.NewString(); return &s }()
	}

	err := chainmwcli.CreateChain(ctx, &chainmwpb.ChainReq{
		EntID:      h.EntID,
		ChainType:  h.ChainType,
		NativeUnit: h.NativeUnit,
		AtomicUnit: h.AtomicUnit,
		UnitExp:    h.UnitExp,
		ENV:        h.ENV,
		ChainID:    h.ChainID,
		NickName:   h.Nickname,
		GasType:    h.GasType,
	})
	if err != nil {
		return nil, err
	}

	return h.GetChain(ctx)
}
