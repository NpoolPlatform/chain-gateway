package tx

import (
	"context"

	txmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/tx"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/tx"
	txmgrpb "github.com/NpoolPlatform/message/npool/chain/mgr/v1/tx"

	accmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/account"
	accmgrpb "github.com/NpoolPlatform/message/npool/account/mgr/v1/account"
)

func GetTxs(ctx context.Context, offset, limit int32) ([]*npool.Tx, uint32, error) {
	infos, total, err := txmwcli.GetTxs(ctx, &txmgrpb.Conds{}, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	ids := []string{}
	for _, info := range infos {
		ids = append(ids, info.FromAccountID, info.ToAccountID)
	}

	accs, _, err := accmwcli.GetManyAccounts(ctx, ids)
	if err != nil {
		return nil, 0, err
	}

	accMap := map[string]*accmgrpb.Account{}
	for _, acc := range accs {
		accMap[acc.ID] = acc
	}

	txs := []*npool.Tx{}
	for _, info := range infos {
		from, ok := accMap[info.FromAccountID]
		if !ok {
			continue
		}

		to, ok := accMap[info.ToAccountID]
		if !ok {
			continue
		}

		txs = append(txs, &npool.Tx{
			ID:            info.ID,
			CoinTypeID:    info.CoinTypeID,
			CoinName:      info.CoinName,
			CoinLogo:      info.CoinLogo,
			CoinUnit:      info.CoinUnit,
			CoinENV:       info.CoinENV,
			FromAccountID: info.FromAccountID,
			FromUsedFor:   from.UsedFor,
			FromAddress:   from.Address,
			ToAccountID:   info.ToAccountID,
			ToUsedFor:     to.UsedFor,
			ToAddress:     to.Address,
			Amount:        info.Amount,
			FeeAmount:     info.FeeAmount,
			ChainTxID:     info.ChainTxID,
			State:         info.State,
			Extra:         info.Extra,
			Type:          info.Type,
			CreatedAt:     info.CreatedAt,
			UpdatedAt:     info.UpdatedAt,
		})
	}

	return txs, total, nil
}
