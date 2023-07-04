package coin

import (
	"context"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"
)

func (h *Handler) UpdateCoin(ctx context.Context) (*coinmwpb.Coin, error) {
	return coinmwcli.UpdateCoin(ctx, &coinmwpb.CoinReq{
		ID:                          h.ID,
		Logo:                        h.Logo,
		Presale:                     h.Presale,
		ReservedAmount:              h.ReservedAmount,
		ForPay:                      h.ForPay,
		HomePage:                    h.HomePage,
		Specs:                       h.Specs,
		FeeCoinTypeID:               h.FeeCoinTypeID,
		WithdrawFeeByStableUSD:      h.WithdrawFeeByStableUSD,
		WithdrawFeeAmount:           h.WithdrawFeeAmount,
		CollectFeeAmount:            h.CollectFeeAmount,
		HotWalletFeeAmount:          h.HotWalletFeeAmount,
		HotWalletAccountAmount:      h.HotWalletAccountAmount,
		LowFeeAmount:                h.LowFeeAmount,
		HotLowFeeAmount:             h.HotLowFeeAmount,
		PaymentAccountCollectAmount: h.PaymentAccountCollectAmount,
		Disabled:                    h.Disabled,
		StableUSD:                   h.StableUSD,
		LeastTransferAmount:         h.LeastTransferAmount,
		NeedMemo:                    h.NeedMemo,
		RefreshCurrency:             h.RefreshCurrency,
		CheckNewAddressBalance:      h.CheckNewAddressBalance,
	})
}
