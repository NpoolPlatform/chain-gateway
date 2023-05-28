package currencyhistory

import (
	"context"

	currencyhismwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency/history"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	currencymwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/currency"
	currencyhismwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/currency/history"
)

func (h *Handler) GetCurrencies(ctx context.Context) ([]*currencymwpb.Currency, uint32, error) {
	conds := &currencyhismwpb.Conds{}
	if len(h.CoinTypeIDs) > 0 {
		conds.CoinTypeIDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.CoinTypeIDs}
	}
	if h.StartAt != nil {
		conds.StartAt = &basetypes.Uint32Val{Op: cruder.LTE, Value: *h.StartAt}
	}
	if h.EndAt != nil {
		conds.EndAt = &basetypes.Uint32Val{Op: cruder.LTE, Value: *h.EndAt}
	}
	return currencyhismwcli.GetCurrencies(ctx, conds, h.Offset, h.Limit)
}