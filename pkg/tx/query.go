package tx

import (
	"context"

	accmwcli "github.com/NpoolPlatform/account-middleware/pkg/client/account"
	txmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/tx"
	wlog "github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	accmwpb "github.com/NpoolPlatform/message/npool/account/mw/v1/account"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/tx"
	txmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/tx"
)

type queryHandler struct {
	*Handler
	infos []*txmwpb.Tx
	total uint32
}

func (h *queryHandler) formalizeAccounts(ctx context.Context, txs []*npool.Tx) ([]*npool.Tx, error) {
	ids := []string{}
	for _, info := range txs {
		ids = append(ids, info.FromAccountID, info.ToAccountID)
	}

	accs, _, err := accmwcli.GetAccounts(ctx, &accmwpb.Conds{
		EntIDs: &basetypes.StringSliceVal{Op: cruder.IN, Value: ids},
	}, 0, int32(len(ids)))
	if err != nil {
		return nil, err
	}

	accMap := map[string]*accmwpb.Account{}
	for _, acc := range accs {
		accMap[acc.EntID] = acc
	}

	for _, info := range txs {
		from, ok := accMap[info.FromAccountID]
		if !ok {
			continue
		}

		to, ok := accMap[info.ToAccountID]
		if !ok {
			continue
		}

		info.FromUsedFor = from.UsedFor
		info.FromAddress = from.Address
		info.ToUsedFor = to.UsedFor
		info.ToAddress = to.Address
	}

	return txs, nil
}

func (h *queryHandler) formalizeApps(ctx context.Context, txs []*npool.Tx) ([]*npool.Tx, error) {
	// TODO: here we have to expand app according to tx type
	return txs, nil
}

func (h *queryHandler) formalize(ctx context.Context) ([]*npool.Tx, error) {
	infos := []*npool.Tx{}
	for _, info := range h.infos {
		infos = append(infos, &npool.Tx{
			ID:            info.ID,
			EntID:         info.EntID,
			CoinTypeID:    info.CoinTypeID,
			CoinName:      info.CoinName,
			CoinLogo:      info.CoinLogo,
			CoinUnit:      info.CoinUnit,
			CoinENV:       info.CoinENV,
			FromAccountID: info.FromAccountID,
			ToAccountID:   info.ToAccountID,
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

	infos, err := h.formalizeAccounts(ctx, infos)
	if err != nil {
		return nil, err
	}
	infos, err = h.formalizeApps(ctx, infos)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (h *Handler) GetTx(ctx context.Context) (*npool.Tx, error) {
	info, err := txmwcli.GetTx(ctx, *h.EntID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, wlog.Errorf("invalid tx")
	}

	handler := &queryHandler{
		Handler: h,
		infos:   []*txmwpb.Tx{info},
	}

	_infos, err := handler.formalize(ctx)
	if err != nil {
		return nil, err
	}

	return _infos[0], nil
}

func (h *Handler) GetTxs(ctx context.Context) ([]*npool.Tx, uint32, error) {
	infos, total, err := txmwcli.GetTxs(ctx, &txmwpb.Conds{}, h.Offset, h.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(infos) == 0 {
		return nil, 0, nil
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
