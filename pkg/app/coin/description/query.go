package description

import (
	"context"

	descmwcli "github.com/NpoolPlatform/chain-middleware/pkg/client/app/coin/description"
	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
	descmwpb "github.com/NpoolPlatform/message/npool/chain/mw/v1/app/coin/description"
)

func (h *Handler) GetCoinDescriptions(ctx context.Context) ([]*descmwpb.CoinDescription, uint32, error) {
	conds := &descmwpb.Conds{}
	if h.AppID != nil {
		conds.AppID = &basetypes.StringVal{Op: cruder.EQ, Value: *h.AppID}
	}
	return descmwcli.GetCoinDescriptions(ctx, conds, h.Offset, h.Limit)
}
