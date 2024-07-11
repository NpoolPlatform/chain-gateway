package tx

import (
	"context"

	txmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/tx"
	wlog "github.com/NpoolPlatform/go-service-framework/pkg/wlog"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	npool "github.com/NpoolPlatform/message/npool/chain/gw/v1/tx"
	txmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/tx"
	sphinxproxypb "github.com/NpoolPlatform/message/npool/sphinxproxy"
	sphinxproxycli "github.com/NpoolPlatform/sphinx-proxy/pkg/client"
)

func (h *Handler) UpdateTx(ctx context.Context) (*npool.Tx, error) {
	tx, err := txmwcli.GetTxOnly(ctx, &txmwpb.Conds{
		// ID:    &basetypes.Uint32Val{Op: cruder.EQ, Value: *h.ID},
		EntID: &basetypes.StringVal{Op: cruder.EQ, Value: *h.EntID},
	})
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, wlog.Errorf("invalid tx")
	}
	if tx.State != basetypes.TxState_TxStateFail {
		return nil, wlog.Errorf("permission denied")
	}

	tx1, err := sphinxproxycli.GetTransaction(ctx, tx.EntID)
	if err != nil {
		return nil, err
	}
	if tx1 == nil {
		return nil, wlog.Errorf("invalid tx")
	}
	if tx1.TransactionState != sphinxproxypb.TransactionState_TransactionStateFail {
		return nil, wlog.Errorf("permission denied")
	}

	if _, err := txmwcli.UpdateTx(ctx, &txmwpb.TxReq{
		ID:    h.ID,
		EntID: h.EntID,
		State: h.State,
	}); err != nil {
		return nil, wlog.WrapError(err)
	}

	if err := sphinxproxycli.UpdateTransaction(ctx, &sphinxproxypb.UpdateTransactionRequest{
		TransactionID:        tx.EntID,
		TransactionState:     tx1.TransactionState,
		NextTransactionState: sphinxproxypb.TransactionState_TransactionStateWait,
	}); err != nil {
		return nil, wlog.WrapError(err)
	}

	return h.GetTx(ctx)
}
