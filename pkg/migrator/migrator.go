//nolint:nolintlint
package migrator

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/chain-middleware/pkg/db"
	"github.com/NpoolPlatform/chain-middleware/pkg/db/ent"
	entcurrencyfeed "github.com/NpoolPlatform/chain-middleware/pkg/db/ent/currencyfeed"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	redis2 "github.com/NpoolPlatform/go-service-framework/pkg/redis"

	servicename "github.com/NpoolPlatform/chain-gateway/pkg/servicename"
	basetypes "github.com/NpoolPlatform/message/npool/basetypes/v1"
)

func lockKey() string {
	const keyServiceID = "serviceid"
	serviceID := config.GetStringValueWithNameSpace(servicename.ServiceDomain, keyServiceID)
	return fmt.Sprintf("migrator:%v:%v", servicename.ServiceDomain, serviceID)
}

//nolint:funlen
func Migrate(ctx context.Context) error {
	if err := db.Init(); err != nil {
		return err
	}

	if err := redis2.TryLock(lockKey(), 0); err != nil {
		return err
	}
	defer func() {
		_ = redis2.Unlock(lockKey())
	}()

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
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
}
