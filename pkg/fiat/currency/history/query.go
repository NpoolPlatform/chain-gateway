package currencyhistory

import (
	"context"

	currencyhismwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/fiat/currency/history"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	currencymwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/fiat/currency"
	currencyhismwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/fiat/currency/history"
)

func (h *Handler) GetCurrencies(ctx context.Context) ([]*currencymwpb.Currency, uint32, error) {
	conds := &currencyhismwpb.Conds{}
	if len(h.FiatIDs) > 0 {
		conds.FiatIDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.FiatIDs}
	}
	if h.StartAt != nil {
		conds.StartAt = &basetypes.Uint32Val{Op: cruder.GTE, Value: *h.StartAt}
	}
	if h.EndAt != nil {
		conds.EndAt = &basetypes.Uint32Val{Op: cruder.LTE, Value: *h.EndAt}
	}
	return currencyhismwcli.GetCurrencies(ctx, conds, h.Offset, h.Limit)
}
