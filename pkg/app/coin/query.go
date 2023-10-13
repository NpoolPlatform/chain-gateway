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
	infos []*appcoinmwpb.Coin
	total uint32
}

func (h *queryHandler) formalize(ctx context.Context) ([]*npool.Coin, error) {
	ids := []string{}
	for _, info := range h.infos {
		ids = append(ids, info.CoinTypeID)
	}

	conds := &appdefaultgoodmwpb.Conds{
		CoinTypeIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	infos, _, err := appdefaultgoodmwcli.GetDefaults(ctx, conds, 0, int32(len(ids)))
	if err != nil {
		return nil, err
	}

	infoMap := map[string]*appdefaultgoodmwpb.Default{}
	for _, info := range infos {
		infoMap[info.AppID+info.CoinTypeID] = info
	}

	_infos := []*npool.Coin{}
	for _, info := range h.infos {
		_info := &npool.Coin{
			ID:                          info.ID,
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

		dinfo, ok := infoMap[info.AppID+info.CoinTypeID]
		if ok {
			_info.DefaultGoodID = &dinfo.AppGoodID
		}

		_infos = append(_infos, _info)
	}
	return _infos, nil
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
		infos:   []*appcoinmwpb.Coin{info},
	}

	infos, err := handler.formalize(ctx)
	if err != nil {
		return nil, err
	}

	return infos[0], nil
}

func (h *Handler) GetCoins(ctx context.Context) ([]*npool.Coin, uint32, error) {
	conds := &appcoinmwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}

	infos, total, err := appcoinmwcli.GetCoins(ctx, conds, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}

	handler := &queryHandler{
		Handler: h,
		infos:   infos,
		total:   total,
	}

	_infos, err := handler.formalize(ctx)
	if err != nil {
		return nil, 0, err
	}

	return _infos, total, nil
}

func (h *Handler) GetCoinExt(ctx context.Context, info *appcoinmwpb.Coin) (*npool.Coin, error) {
	h.AppID = &info.AppID

	handler := &queryHandler{
		Handler: h,
		infos:   []*appcoinmwpb.Coin{info},
	}

	infos, err := handler.formalize(ctx)
	if err != nil {
		return nil, err
	}

	return infos[0], nil
}
