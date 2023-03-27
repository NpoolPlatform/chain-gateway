//nolint:nolintlint
package migrator

import (
	"context"

	"github.com/NpoolPlatform/chain-manager/pkg/db"
	"github.com/NpoolPlatform/chain-manager/pkg/db/ent"
	entsetting "github.com/NpoolPlatform/chain-manager/pkg/db/ent/setting"

	"github.com/shopspring/decimal"
)

func Migrate(ctx context.Context) error {
	if err := db.Init(); err != nil {
		return err
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		_, err := tx.ExecContext(
			_ctx,
			"update settings set hot_low_fee_amount='0' where hot_low_fee_amount is NULL",
		)
		if err != nil {
			return err
		}

		infos, err := tx.
			Setting.
			Query().
			Where(
				entsetting.DeletedAt(0),
			).
			All(_ctx)
		if err != nil {
			return err
		}

		for _, info := range infos {
			if info.HotLowFeeAmount.Cmp(decimal.NewFromInt(0)) > 0 {
				continue
			}

			_, err := tx.
				Setting.
				UpdateOneID(info.ID).
				SetHotLowFeeAmount(info.LowFeeAmount.Mul(decimal.NewFromInt(3))).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
