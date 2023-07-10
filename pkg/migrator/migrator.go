//nolint:nolintlint
package migrator

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/chain-middleware/pkg/db"
	"github.com/NpoolPlatform/chain-middleware/pkg/db/ent"
	entcoinbase "github.com/NpoolPlatform/chain-middleware/pkg/db/ent/coinbase"
	entsetting "github.com/NpoolPlatform/chain-middleware/pkg/db/ent/setting"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	servicename "github.com/NpoolPlatform/chain-gateway/pkg/servicename"

	"github.com/google/uuid"
)

func lockKey() string {
	const keyServiceID = "serviceid"
	serviceID := config.GetStringValueWithNameSpace(servicename.ServiceDomain, keyServiceID)
	return fmt.Sprintf("migrator:%v:%v", servicename.ServiceDomain, serviceID)
}

func migrateSetting(ctx context.Context) error {
	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		coinName := "ironfish"
		coins, err := tx.
			CoinBase.
			Query().
			Where(
				entcoinbase.NameContains(coinName),
			).
			All(_ctx)
		if err != nil {
			return err
		}

		ids := []uuid.UUID{}
		for _, coin := range coins {
			ids = append(ids, coin.ID)
		}

		boolTrue := true
		boolFalse := false
		_, err = tx.
			Setting.
			Update().
			Where(
				entsetting.CoinTypeIDIn(ids...),
				entsetting.CheckNewAddressBalance(boolTrue),
			).
			SetCheckNewAddressBalance(boolFalse).
			Save(_ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

//nolint:funlen,gocyclo
func Migrate(ctx context.Context) (err error) {
	logger.Sugar().Infow(
		"Migrate",
		"State", "Start...",
		"LockKey", lockKey(),
	)
	defer logger.Sugar().Infow(
		"Migrate",
		"State", "Done...",
		"Error", err,
	)

	if err = db.Init(); err != nil { //nolint
		return err
	}

	for {
		if err = redis2.TryLock(lockKey(), 0); err != nil { //nolint
			logger.Sugar().Infow(
				"Migrate",
				"State", "Waiting...",
				"Error", err,
			)
			time.Sleep(time.Minute)
			continue
		}
		break
	}
	defer func() {
		_ = redis2.Unlock(lockKey())
	}()

	if err := migrateSetting(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	logger.Sugar().Infow("Migrate", "Done", "success")

	return nil
}
