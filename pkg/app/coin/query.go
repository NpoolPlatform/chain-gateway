package appcoin

import (
	"context"
	"fmt"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	appdefaultgoodmwcli "github.com/NpoolPlatform/good-middleware/pkg/client/app/good/default"
	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/app/coin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin"
	appdefaultgoodmwpb "github.com/NpoolPlatform/message/npool/good/mw/v1/app/good/default"
)

type queryHandler struct {
	*Handler
	coins        []*appcoinmwpb.Coin
	infos        []*npool.Coin
	defaultGoods map[string]*appdefaultgoodmwpb.Default
	total        uint32
}

func (h *queryHandler) getDefaultGoods(ctx context.Context) error {
	coinTypeIDs := func() (_coinTypeIDs []string) {
		for _, info := range h.coins {
			_coinTypeIDs = append(_coinTypeIDs, info.CoinTypeID)
		}
		return
	}()
	conds := &appdefaultgoodmwpb.Conds{
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: coinTypeIDs},
	}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	defaultGoods, _, err := appdefaultgoodmwcli.GetDefaults(ctx, conds, 0, int32(len(coinTypeIDs)))
	if err != nil {
		return err
	}
	h.defaultGoods = map[string]*appdefaultgoodmwpb.Default{}
	for _, defaultGood := range defaultGoods {
		h.defaultGoods[defaultGood.AppID+defaultGood.CoinTypeID] = defaultGood
	}
	return nil
}

func (h *queryHandler) formalize() {
	for _, info := range h.coins {
		_info := &npool.Coin{
			ID:                          info.ID,
			EntID:                       info.EntID,
			AppID:                       info.AppID,
			CoinTypeID:                  info.CoinTypeID,
			Name:                        info.Name,
			CoinName:                    info.CoinName,
			DisplayNames:                info.DisplayNames,
			Logo:                        info.Logo,
			Unit:                        info.Unit,
			Presale:                     info.Presale,
			ReservedAmount:              info.ReservedAmount,
			ForPay:                      info.ForPay,
			ProductPage:                 info.ProductPage,
			CoinForPay:                  info.CoinForPay,
			ENV:                         info.ENV,
			HomePage:                    info.HomePage,
			Specs:                       info.Specs,
			StableUSD:                   info.StableUSD,
			FeeCoinTypeID:               info.FeeCoinTypeID,
			FeeCoinName:                 info.FeeCoinName,
			FeeCoinLogo:                 info.FeeCoinLogo,
			FeeCoinUnit:                 info.FeeCoinUnit,
			FeeCoinENV:                  info.FeeCoinENV,
			WithdrawFeeByStableUSD:      info.WithdrawFeeByStableUSD,
			WithdrawFeeAmount:           info.WithdrawFeeAmount,
			CollectFeeAmount:            info.CollectFeeAmount,
			HotWalletFeeAmount:          info.HotWalletFeeAmount,
			LowFeeAmount:                info.LowFeeAmount,
			HotWalletAccountAmount:      info.HotWalletAccountAmount,
			PaymentAccountCollectAmount: info.PaymentAccountCollectAmount,
			WithdrawAutoReviewAmount:    info.WithdrawAutoReviewAmount,
			MarketValue:                 info.MarketValue,
			SettleValue:                 info.SettleValue,
			SettlePercent:               info.SettlePercent,
			SettleTipsStr:               info.SettleTipsStr,
			SettleTips:                  info.SettleTips,
			Setter:                      info.Setter,
			Disabled:                    info.Disabled,
			CoinDisabled:                info.CoinDisabled,
			CreatedAt:                   info.CreatedAt,
			UpdatedAt:                   info.UpdatedAt,
			Display:                     info.Display,
			DisplayIndex:                info.DisplayIndex,
			MaxAmountPerWithdraw:        info.MaxAmountPerWithdraw,
			LeastTransferAmount:         info.LeastTransferAmount,
			NeedMemo:                    info.NeedMemo,
		}
		defaultGood, ok := h.defaultGoods[info.AppID+info.CoinTypeID]
		if ok {
			_info.DefaultGoodID = &defaultGood.AppGoodID
		}
		h.infos = append(h.infos, _info)
	}
}

func (h *Handler) GetCoin(ctx context.Context) (*npool.Coin, error) {
	if h.EntID == nil {
		return nil, fmt.Errorf("invalid entid")
	}

	info, err := appcoinmwcli.GetCoin(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	h.AppID = &info.AppID

	handler := &queryHandler{
		Handler: h,
		coins:   []*appcoinmwpb.Coin{info},
	}

	if err := handler.getDefaultGoods(ctx); err != nil {
		return nil, err
	}

	handler.formalize()

	return handler.infos[0], nil
}

func (h *Handler) GetCoins(ctx context.Context) ([]*npool.Coin, uint32, error) {
	conds := &appcoinmwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	if h.ForPay != nil {
		conds.ForPay = &basetypes.BoolVal{Op: cruder.EQ, Value: *h.ForPay}
	}

	coins, total, err := appcoinmwcli.GetCoins(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}

	handler := &queryHandler{
		Handler: h,
		coins:   coins,
		total:   total,
	}

	if err := handler.getDefaultGoods(ctx); err != nil {
		return nil, 0, err
	}

	handler.formalize()

	return handler.infos, total, nil
}

func (h *Handler) GetCoinExt(ctx context.Context, info *appcoinmwpb.Coin) (*npool.Coin, error) {
	h.AppID = &info.AppID

	handler := &queryHandler{
		Handler: h,
		coins:   []*appcoinmwpb.Coin{info},
	}
	if err := handler.getDefaultGoods(ctx); err != nil {
		return nil, err
	}

	handler.formalize()

	return handler.infos[0], nil
}
