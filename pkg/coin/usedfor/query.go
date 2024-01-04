package coinusedfor

import (
	"context"

	coinusedformwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/usedfor"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	coinusedformwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/usedfor"
)

func (h *Handler) GetCoinUsedFors(ctx context.Context) ([]*coinusedformwpb.CoinUsedFor, uint32, error) {
	conds := &coinusedformwpb.Conds{}
	if h.CoinTypeID != nil {
		conds.CoinTypeID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinTypeID}
	}
	if len(h.CoinTypeIDs) > 0 {
		conds.CoinTypeIDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.CoinTypeIDs}
	}
	if h.UsedFor != nil {
		conds.UsedFor = &basetypes.Uint32Val{Op: cruder.EQ, Value: uint32(*h.UsedFor)}
	}
	if len(h.UsedFors) > 0 {
		_usedFors := []uint32{}
		for _, usedFor := range h.UsedFors {
			_usedFors = append(_usedFors, uint32(usedFor))
		}
		conds.UsedFors = &basetypes.Uint32SliceVal{Op: cruder.IN, Value: _usedFors}
	}
	return coinusedformwcli.GetCoinUsedFors(ctx, conds, h.Offset, h.Limit)
}
