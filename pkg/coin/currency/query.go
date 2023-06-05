package currency

import (
	"context"

	currencymwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	currencymwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin/currency"
)

func (h *Handler) GetCurrency(ctx context.Context) (*currencymwpb.Currency, error) {
	return currencymwcli.GetCurrencyOnly(ctx, &currencymwpb.Conds{
		CoinTypeID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.CoinTypeID},
	})
}

func (h *Handler) GetCurrencies(ctx context.Context) ([]*currencymwpb.Currency, uint32, error) {
	conds := &currencymwpb.Conds{}
	if len(h.CoinTypeIDs) > 0 {
		conds.CoinTypeIDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.CoinTypeIDs}
	}
	return currencymwcli.GetCurrencies(ctx, conds, h.Offset, h.Limit)
}
