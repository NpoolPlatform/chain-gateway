package tx

import (
	"context"

	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/tx"
)

func GetTxs(ctx context.Context, appID string, offset, limit int32) ([]*npool.Tx, uint32, error) {
	return nil, 0, nil
}
