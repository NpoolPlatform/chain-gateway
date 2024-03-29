package appcoin

import (
	"context"
	"fmt"

	appcoinmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/app/coin"
)

func (h *Handler) DeleteCoin(ctx context.Context) (*npool.Coin, error) {
	if h.ID == nil {
		return nil, fmt.Errorf("invalid id")
	}

	// TODO: check appid / cointypeid / id

	deleteInfo, err := appcoinmwcli.DeleteCoin(ctx, *h.ID)
	if err != nil {
		return nil, err
	}

	if deleteInfo == nil {
		return nil, nil
	}

	info, err := h.GetCoinExt(ctx, deleteInfo)
	if err != nil {
		return nil, err
	}

	return info, nil
}
