//nolint:dupl
package appcoin

import (
	"context"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/appcoin"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	"github.com/NpoolPlatform/good-middleware/pkg/client/appdefaultgood"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/appcoin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/appcoin"
	appgoodmgrpb "github.com/NpoolPlatform/message/npool/good/mgr/v1/appdefaultgood"
)

func GetAppCoin(ctx context.Context, id string) (*npool.Coin, error) {
	row, err := appcoinmwcli.GetCoin(ctx, id)
	if err != nil {
		logger.Sugar().Errorw("GetAppCoins", "error", err)
		return nil, err
	}

	goodInfo, err := appdefaultgood.GetAppDefaultGoodOnly(ctx, &appgoodmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: row.AppID,
		},
		CoinTypeID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: row.CoinTypeID,
		},
	})
	if err != nil {
		return nil, err
	}
	goodID := ""
	if goodInfo != nil {
		goodID = goodInfo.GoodID
	}
	return &npool.Coin{
		ID:                          row.ID,
		AppID:                       row.AppID,
		CoinTypeID:                  row.CoinTypeID,
		Name:                        row.Name,
		CoinName:                    row.CoinName,
		DisplayNames:                row.DisplayNames,
		Logo:                        row.Logo,
		Unit:                        row.Unit,
		Presale:                     row.Presale,
		ReservedAmount:              row.ReservedAmount,
		ForPay:                      row.ForPay,
		ProductPage:                 row.ProductPage,
		CoinForPay:                  row.CoinForPay,
		ENV:                         row.ENV,
		HomePage:                    row.HomePage,
		Specs:                       row.Specs,
		StableUSD:                   row.StableUSD,
		FeeCoinTypeID:               row.FeeCoinTypeID,
		FeeCoinName:                 row.FeeCoinName,
		FeeCoinLogo:                 row.FeeCoinLogo,
		FeeCoinUnit:                 row.FeeCoinUnit,
		FeeCoinENV:                  row.FeeCoinENV,
		WithdrawFeeByStableUSD:      row.WithdrawFeeByStableUSD,
		WithdrawFeeAmount:           row.WithdrawFeeAmount,
		CollectFeeAmount:            row.CollectFeeAmount,
		HotWalletFeeAmount:          row.HotWalletFeeAmount,
		LowFeeAmount:                row.LowFeeAmount,
		HotWalletAccountAmount:      row.HotWalletAccountAmount,
		PaymentAccountCollectAmount: row.PaymentAccountCollectAmount,
		WithdrawAutoReviewAmount:    row.WithdrawAutoReviewAmount,
		MarketValue:                 row.MarketValue,
		SettleValue:                 row.SettleValue,
		SettlePercent:               row.SettlePercent,
		SettleTipsStr:               row.SettleTipsStr,
		SettleTips:                  row.SettleTips,
		Setter:                      row.Setter,
		Disabled:                    row.Disabled,
		CoinDisabled:                row.CoinDisabled,
		CreatedAt:                   row.CreatedAt,
		UpdatedAt:                   row.UpdatedAt,
		DailyRewardAmount:           row.DailyRewardAmount,
		Display:                     row.Display,
		DisplayIndex:                row.DisplayIndex,
		MaxAmountPerWithdraw:        row.MaxAmountPerWithdraw,
		LeastTransferAmount:         row.LeastTransferAmount,
		DefaultGoodID:               goodID,
		NeedMemo:                    row.NeedMemo,
	}, nil
}
