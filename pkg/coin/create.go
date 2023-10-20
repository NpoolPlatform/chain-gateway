package coin

import (
	"context"

	coinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/coin"
	coinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/coin"
)

func (h *Handler) CreateCoin(ctx context.Context) (*coinmwpb.Coin, error) {
	return coinmwcli.CreateCoin(ctx, &coinmwpb.CoinReq{
		Name:                h.Name,
		Unit:                h.Unit,
		ENV:                 h.ENV,
		ChainType:           h.ChainType,
		ChainNativeUnit:     h.ChainNativeUnit,
		ChainAtomicUnit:     h.ChainAtomicUnit,
		ChainUnitExp:        h.ChainUnitExp,
		GasType:             h.GasType,
		ChainID:             h.ChainID,
		ChainNickname:       h.ChainNickname,
		ChainNativeCoinName: h.ChainNativeCoinName,
	})
}
