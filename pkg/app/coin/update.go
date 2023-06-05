package appcoin

import (
	"context"
	"fmt"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/app/coin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
)

func (h *Handler) UpdateCoin(ctx context.Context) (*npool.Coin, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	// TODO: check appid / cointypeid / id

	_, err := appcoinmwcli.UpdateCoin(ctx, &appcoinmwpb.CoinReq{
		ID:                       h.ID,
		Name:                     h.Name,
		DisplayNames:             h.DisplayNames,
		Logo:                     h.Logo,
		ForPay:                   h.ForPay,
		WithdrawAutoReviewAmount: h.WithdrawAutoReviewAmount,
		MarketValue:              h.MarketValue,
		SettlePercent:            h.SettlePercent,
		SettleTips:               h.SettleTips,
		DailyRewardAmount:        h.DailyRewardAmount,
		ProductPage:              h.ProductPage,
		Disabled:                 h.Disabled,
		Display:                  h.Display,
		DisplayIndex:             h.DisplayIndex,
		MaxAmountPerWithdraw:     h.MaxAmountPerWithdraw,
	})
	if err != nil {
		return nil, err
	}

	return h.GetCoin(ctx)
}
