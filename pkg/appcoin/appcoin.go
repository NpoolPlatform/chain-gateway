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

func GetAppCoins(ctx context.Context, appID string, offset, limit int32) ([]*npool.Coin, uint32, error) {
	rows, total, err := appcoinmwcli.GetCoins(ctx, &appcoinmwpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
	}, offset, limit)
	if err != nil {
		logger.Sugar().Errorw("GetAppCoins", "error", err)
		return nil, 0, err
	}
	coinTypeIDs := []string{}
	for _, val := range rows {
		coinTypeIDs = append(coinTypeIDs, val.GetCoinTypeID())
	}
	goodInfos, _, err := appdefaultgood.GetAppDefaultGoods(ctx, &appgoodmgrpb.Conds{
		AppID: &commonpb.StringVal{
			Op:    cruder.EQ,
			Value: appID,
		},
		CoinTypeIDs: &commonpb.StringSliceVal{
			Op:    cruder.IN,
			Value: coinTypeIDs,
		},
	}, 0, int32(len(rows)))
	if err != nil {
		return nil, 0, err
	}

	defaultGoodMap := map[string]*appgoodmgrpb.AppDefaultGood{}
	for _, val := range goodInfos {
		defaultGoodMap[val.CoinTypeID] = val
	}
	infos := []*npool.Coin{}
	for _, val := range rows {
		goodID := ""
		defaultGood, ok := defaultGoodMap[val.CoinTypeID]
		if ok {
			goodID = defaultGood.GoodID
		}
		infos = append(infos, &npool.Coin{
			ID:                          val.ID,
			AppID:                       val.AppID,
			CoinTypeID:                  val.CoinTypeID,
			Name:                        val.Name,
			CoinName:                    val.CoinName,
			DisplayNames:                val.DisplayNames,
			Logo:                        val.Logo,
			Unit:                        val.Unit,
			Presale:                     val.Presale,
			ReservedAmount:              val.ReservedAmount,
			ForPay:                      val.ForPay,
			ProductPage:                 val.ProductPage,
			CoinForPay:                  val.CoinForPay,
			ENV:                         val.ENV,
			HomePage:                    val.HomePage,
			Specs:                       val.Specs,
			StableUSD:                   val.StableUSD,
			FeeCoinTypeID:               val.FeeCoinTypeID,
			FeeCoinName:                 val.FeeCoinName,
			FeeCoinLogo:                 val.FeeCoinLogo,
			FeeCoinUnit:                 val.FeeCoinUnit,
			FeeCoinENV:                  val.FeeCoinENV,
			WithdrawFeeByStableUSD:      val.WithdrawFeeByStableUSD,
			WithdrawFeeAmount:           val.WithdrawFeeAmount,
			CollectFeeAmount:            val.CollectFeeAmount,
			HotWalletFeeAmount:          val.HotWalletFeeAmount,
			LowFeeAmount:                val.LowFeeAmount,
			HotWalletAccountAmount:      val.HotWalletAccountAmount,
			PaymentAccountCollectAmount: val.PaymentAccountCollectAmount,
			WithdrawAutoReviewAmount:    val.WithdrawAutoReviewAmount,
			MarketValue:                 val.MarketValue,
			SettleValue:                 val.SettleValue,
			SettlePercent:               val.SettlePercent,
			SettleTipsStr:               val.SettleTipsStr,
			SettleTips:                  val.SettleTips,
			Setter:                      val.Setter,
			Disabled:                    val.Disabled,
			CoinDisabled:                val.CoinDisabled,
			CreatedAt:                   val.CreatedAt,
			UpdatedAt:                   val.UpdatedAt,
			DailyRewardAmount:           val.DailyRewardAmount,
			Display:                     val.Display,
			DisplayIndex:                val.DisplayIndex,
			MaxAmountPerWithdraw:        val.MaxAmountPerWithdraw,
			LeastTransferAmount:         val.LeastTransferAmount,
			DefaultGoodID:               goodID,
		})
	}
	return infos, total, nil
}

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
	}, nil
}

func CreateAppCoin(ctx context.Context, appID, coinTypeID string) (*npool.Coin, error) {
	row, err := appcoinmwcli.CreateCoin(ctx, &appcoinmwpb.CoinReq{
		AppID:      &appID,
		CoinTypeID: &coinTypeID,
	})
	if err != nil {
		return nil, err
	}
	info, err := GetAppCoin(ctx, row.ID)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func UpdateAppCoin(ctx context.Context, coinReq *appcoinmwpb.CoinReq) (*npool.Coin, error) {
	row, err := appcoinmwcli.UpdateCoin(ctx, coinReq)
	if err != nil {
		return nil, err
	}
	info, err := GetAppCoin(ctx, row.ID)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func DeleteAppCoin(ctx context.Context, id string) (*npool.Coin, error) {
	info, err := GetAppCoin(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = appcoinmwcli.DeleteCoin(ctx, id)
	if err != nil {
		return nil, err
	}
	return info, nil
}
