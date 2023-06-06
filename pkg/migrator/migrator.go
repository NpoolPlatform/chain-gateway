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

	if err = db.Init(); err != nil {
		return err
	}

	if err = redis2.TryLock(lockKey(), 0); err != nil {
		return err
	}
	defer func() {
		_ = redis2.Unlock(lockKey())
	}()

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

	limit := 1000
	err = db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		kept := map[string]uuid.UUID{}

		for {
			logger.Sugar().Errorw(
				"Migrate",
				"Limit", limit,
				"State", "Query start",
			)

			currencies, err := tx.
				Currency.
				Query().
				Where(
					entcurrency.DeletedAt(0),
				).
				Limit(limit).
				Order(ent.Desc(entcurrency.FieldUpdatedAt)).
				All(_ctx)
			logger.Sugar().Errorw(
				"Migrate",
				"Limit", limit,
				"State", "Query done",
			)
			if err != nil {
				logger.Sugar().Errorw(
					"Migrate",
					"Limit", limit,
					"Error", err,
				)
				return err
			}

			updatable := false
			for _, currency := range currencies {
				keptKey := fmt.Sprintf("%v:%v", currency.CoinTypeID, currency.FeedType)
				_kept, ok := kept[keptKey]
				if !ok || _kept != currency.ID {
					updatable = true
					break
				}
			}

			if !updatable {
				logger.Sugar().Errorw(
					"Migrate",
					"Limit", limit,
					"State", "Done",
				)
				break
			}

			_currencies := []*ent.CurrencyHistoryCreate{}
			for _, currency := range currencies {
				_currencies = append(
					_currencies,
					tx.
						CurrencyHistory.
						Create().
						SetCoinTypeID(currency.CoinTypeID).
						SetFeedType(currency.FeedType).
						SetMarketValueHigh(currency.MarketValueHigh).
						SetMarketValueLow(currency.MarketValueLow).
						SetCreatedAt(currency.CreatedAt).
						SetUpdatedAt(currency.UpdatedAt),
				)
			}

			err = tx.
				CurrencyHistory.
				CreateBulk(_currencies...).
				Exec(_ctx)
			if err != nil {
				return err
			}

			for _, currency := range currencies {
				keptKey := fmt.Sprintf("%v:%v", currency.CoinTypeID, currency.FeedType)
				_kept, ok := kept[keptKey]
				if !ok {
					kept[keptKey] = currency.ID
					continue
				}
				if _kept == currency.ID {
					continue
				}

				_, err = tx.
					Currency.
					UpdateOneID(currency.ID).
					SetDeletedAt(uint32(time.Now().Unix())).
					Save(_ctx)
				if err != nil {
					return err
				}
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

func Abort() {
	logger.Sugar().Errorw(
		"Migrate",
		"State", "Aborted",
	)
	_ = redis2.Unlock(lockKey())
}
