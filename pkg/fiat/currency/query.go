package currency

import (
	"context"

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
