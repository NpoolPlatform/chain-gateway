package coinfiat

import (
	"context"

	coinfiatmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/fiat"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	coinfiatmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/fiat"
)

func (h *Handler) GetCoinFiats(ctx context.Context) ([]*coinfiatmwpb.CoinFiat, uint32, error) {
	conds := &coinfiatmwpb.Conds{}
	if len(h.CoinTypeIDs) > 0 {
		conds.CoinTypeIDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.CoinTypeIDs}
	}
	return coinfiatmwcli.GetCoinFiats(ctx, conds, h.Offset, h.Limit)
}
