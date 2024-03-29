package currency

import (
	"context"
	"fmt"

	currencymwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/fiat/currency"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	currencymwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/fiat/currency"
)

func (h *Handler) GetCurrencies(ctx context.Context) ([]*currencymwpb.Currency, uint32, error) {
	conds := &currencymwpb.Conds{}
	if len(h.FiatIDs) > 0 {
		conds.FiatIDs = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.FiatIDs}
	}
	return currencymwcli.GetCurrencies(ctx, conds, h.Offset, h.Limit)
}

func (h *Handler) GetCurrency(ctx context.Context) (*currencymwpb.Currency, error) {
	if h.FiatName == nil {
		return nil, fmt.Errorf("invalid fiatname")
	}
	conds := &currencymwpb.Conds{
		FiatName: &basetypes.StringVal{Op: cruder.EQ, Value: *h.FiatName},
	}
	return currencymwcli.GetCurrencyOnly(ctx, conds)
}
