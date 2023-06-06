//nolint:nolintlint
package migrator

import (
	"context"
	"fmt"
	"time"

	"github.com/NpoolPlatform/chain-middleware/pkg/db"
	"github.com/NpoolPlatform/chain-middleware/pkg/db/ent"
	entcurrency "github.com/NpoolPlatform/chain-middleware/pkg/db/ent/currency"
	entcurrencyfeed "github.com/NpoolPlatform/chain-middleware/pkg/db/ent/currencyfeed"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	servicename "github.com/NpoolPlatform/chain-gateway/pkg/servicename"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/google/uuid"
)

func lockKey() string {
	const keyServiceID = "serviceid"
	serviceID := config.GetStringValueWithNameSpace(servicename.ServiceDomain, keyServiceID)
	return fmt.Sprintf("migrator:%v:%v", servicename.ServiceDomain, serviceID)
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

	err = db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		_, err := tx.
			ExecContext(
				ctx,
				"SET global sql_mode=(SELECT REPLACE(@@sql_mode, 'ONLY_FULL_GROUP_BY', ''))",
			)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		coins, err := tx.
			CoinBase.
			Query().
			All(_ctx)
		if err != nil {
			return err
		}

		_coinMaps := map[basetypes.CurrencyFeedType]map[string]string{
			basetypes.CurrencyFeedType_CoinBase: {
				"fil":        "FIL",
				"filecoin":   "FIL",
				"tfilecoin":  "FIL",
				"btc":        "BTC",
				"bitcoin":    "BTC",
				"tbitcoin":   "BTC",
				"tethereum":  "ETH",
				"eth":        "ETH",
				"ethereum":   "ETH",
				"sol":        "SOL",
				"solana":     "SOL",
				"tsolana":    "SOL",
				"tusdcerc20": "USDT",
				"usdcerc20":  "USDT",
			},
			basetypes.CurrencyFeedType_CoinGecko: {
				"fil":          "filecoin",
				"filecoin":     "filecoin",
				"tfilecoin":    "filecoin",
				"btc":          "bitcoin",
				"bitcoin":      "bitcoin",
				"tbitcoin":     "bitcoin",
				"tethereum":    "ethereum",
				"eth":          "ethereum",
				"ethereum":     "ethereum",
				"tusdt":        "tether",
				"usdt":         "tether",
				"tusdterc20":   "tether",
				"usdterc20":    "tether",
				"tusdttrc20":   "tether",
				"usdttrc20":    "tether",
				"sol":          "solana",
				"solana":       "solana",
				"tsolana":      "solana",
				"tbinancecoin": "binancecoin",
				"binancecoin":  "binancecoin",
				"tbinanceusd":  "binance-usd",
				"binanceusd":   "binance-usd",
				"ttron":        "tron",
				"tron":         "tron",
				"tusdcerc20":   "usd-coin",
				"usdcerc20":    "usd-coin",
			},
		}

		for _, _coin := range coins {
			for _feedType, _coinMap := range _coinMaps {
				_feedCoinName, ok := _coinMap[_coin.Name]
				if !ok {
					continue
				}

				info, err := tx.
					CurrencyFeed.
					Query().
					Where(
						entcurrencyfeed.CoinTypeID(_coin.ID),
						entcurrencyfeed.FeedType(_feedType.String()),
						entcurrencyfeed.DeletedAt(0),
					).
					Only(_ctx)
				if err != nil {
					if !ent.IsNotFound(err) {
						logger.Sugar().Errorw(
							"Migrate",
							"CoinName", _coin.Name,
							"FeedCoinName", _feedCoinName,
							"CoinTypeID", _coin.ID,
						)
						return err
					}
				}
				if info != nil {
					continue
				}

				if _, err := tx.
					CurrencyFeed.
					Create().
					SetCoinTypeID(_coin.ID).
					SetFeedType(_feedType.String()).
					SetFeedCoinName(_feedCoinName).
					Save(_ctx); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = db.WithClient(ctx, func(_ctx context.Context, cli *ent.Client) error {
		_, err := cli.
			Currency.
			Update().
			SetDeletedAt(uint32(time.Now().Unix())).
			Save(_ctx)
		if err != nil {
			return err
		}

		ids, err := cli.
			Currency.
			Query().
			GroupBy(entcurrency.FieldCoinTypeID).
			Strings(_ctx)
		if err != nil {
			return err
		}

		for _, id := range ids {
			info, err := cli.
				Currency.
				Query().
				Where(
					entcurrency.CoinTypeID(uuid.MustParse(id)),
				).
				Order(
					ent.Desc(entcurrency.FieldUpdatedAt),
				).
				Limit(1).
				Only(_ctx)
			if err != nil {
				return err
			}

			_, err = cli.
				Currency.
				UpdateOneID(info.ID).
				SetDeletedAt(0).
				Save(_ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		_, err := tx.
			FiatCurrency.
			Update().
			SetDeletedAt(uint32(time.Now().Unix())).
			Save(_ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
