package currencyhistory

import (
	"context"
	"time"

	currencyhismwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin/currency/history"
	timedef "github.com/NpoolPlatform/go-service-framework/pkg/const/time"
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
	if len(h.CoinNames) > 0 {
		conds.CoinNames = &basetypes.StringSliceVal{Op: cruder.IN, Value: h.CoinNames}
	}
	startAt := uint32(time.Now().Unix()) - timedef.SecondsPerDay*30 //nolint
	if h.StartAt != nil {
		startAt = *h.StartAt
	}
	conds.StartAt = &basetypes.Uint32Val{Op: cruder.GTE, Value: startAt}
	endAt := uint32(time.Now().Unix())
	if h.EndAt != nil {
		endAt = *h.EndAt
	}
	conds.EndAt = &basetypes.Uint32Val{Op: cruder.LTE, Value: endAt}
	return currencyhismwcli.GetCurrencies(ctx, conds, h.Offset, h.Limit)
}
