package appcoin

import (
	"context"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/app/coin"
	appcoinmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/appcoin"
)

func (h *Handler) CreateCoin(ctx context.Context) (*npool.Coin, error) {
	info, err := appcoinmwcli.CreateCoin(ctx, &appcoinmwpb.CoinReq{
		AppID:      h.AppID,
		CoinTypeID: h.CoinTypeID,
	})
	if err != nil {
		return nil, err
	}

	h.ID = &info.ID

	return h.GetCoin(ctx)
}
