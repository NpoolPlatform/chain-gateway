//nolint:nolintlint
package migrator

import (
	"context"

	"github.com/NpoolPlatform/chain-manager/pkg/db"
	"github.com/NpoolPlatform/chain-manager/pkg/db/ent"

	"github.com/shopspring/decimal"
)

func Migrate(ctx context.Context) error {
	if err := db.Init(); err != nil {
		return err
	}

	return db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		_, err := cli.
			AppCoin.
			Update().
			SetMaxAmountPerWithdraw(decimal.NewFromInt(40000)).
			SetDailyRewardAmount(decimal.NewFromInt(0)).
			Save(_ctx)
		if err != nil {
			return err
		}

		_, err = cli.
			Setting.
			Update().
			SetLeastTransferAmount(decimal.RequireFromString("0.1")).
			Save(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
}
