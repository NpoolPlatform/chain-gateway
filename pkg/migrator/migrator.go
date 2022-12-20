//nolint
package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/NpoolPlatform/chain-manager/pkg/db"
	"github.com/NpoolPlatform/chain-manager/pkg/db/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/shopspring/decimal"

	billingent "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent"
	entcointx "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent/coinaccounttransaction"

	gasfeederent "github.com/NpoolPlatform/gas-feeder/pkg/db/ent"
	oracleent "github.com/NpoolPlatform/oracle-manager/pkg/db/ent"
	projinfoent "github.com/NpoolPlatform/project-info-manager/pkg/db/ent"
	coininfoent "github.com/NpoolPlatform/sphinx-coininfo/pkg/db/ent"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	gasfeederconst "github.com/NpoolPlatform/gas-feeder/pkg/message/const"
	oracleconst "github.com/NpoolPlatform/oracle-manager/pkg/message/const"
	projinfoconst "github.com/NpoolPlatform/project-info-manager/pkg/message/const"
	coininfoconst "github.com/NpoolPlatform/sphinx-coininfo/pkg/message/const"

	appmwcli "github.com/NpoolPlatform/appuser-middleware/pkg/client/app"
	appmwpb "github.com/NpoolPlatform/message/npool/appuser/mw/v1/app"

	descmgrpb "github.com/NpoolPlatform/message/npool/chain/mgr/v1/appcoin/description"
	txmgrpb "github.com/NpoolPlatform/message/npool/chain/mgr/v1/tx"

	_ "github.com/NpoolPlatform/gas-feeder/pkg/db/ent/runtime"
	_ "github.com/NpoolPlatform/oracle-manager/pkg/db/ent/runtime"
	_ "github.com/NpoolPlatform/project-info-manager/pkg/db/ent/runtime"
	_ "github.com/NpoolPlatform/sphinx-coininfo/pkg/db/ent/runtime"

	"github.com/google/uuid"
)

const (
	keyUsername = "username"
	keyPassword = "password"
	keyDBName   = "database_name"
	maxOpen     = 10
	maxIdle     = 10
	MaxLife     = 3
)

func dsn(hostname string) (string, error) {
	username := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyUsername)
	password := config.GetStringValueWithNameSpace(constant.MysqlServiceName, keyPassword)
	dbname := config.GetStringValueWithNameSpace(hostname, keyDBName)

	svc, err := config.PeekService(constant.MysqlServiceName)
	if err != nil {
		logger.Sugar().Warnw("dsb", "error", err)
		return "", err
	}

	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&interpolateParams=true",
		username, password,
		svc.Address,
		svc.Port,
		dbname,
	), nil
}

func open(hostname string) (conn *sql.DB, err error) {
	hdsn, err := dsn(hostname)
	if err != nil {
		return nil, err
	}

	logger.Sugar().Infow("open", "hdsn", hdsn)

	conn, err = sql.Open("mysql", hdsn)
	if err != nil {
		return nil, err
	}

	// https://github.com/go-sql-driver/mysql
	// See "Important settings" section.

	conn.SetConnMaxLifetime(time.Minute * MaxLife)
	conn.SetMaxOpenConns(maxOpen)
	conn.SetMaxIdleConns(maxIdle)

	return conn, nil
}

func migrateTx(ctx context.Context, conn *sql.DB) error {
	cli1 := billingent.NewClient(billingent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	txs, err := cli1.
		CoinAccountTransaction.
		Query().
		Where(
			entcointx.DeleteAt(0),
		).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateTx", "error", err)
		return err
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		infos, err := tx.
			Tran.
			Query().
			Limit(1).
			All(ctx)
		if err != nil {
			logger.Sugar().Errorw("migrateTx", "error", err)
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, tran := range txs {
			found := false
			for _, coin := range coinInfos {
				if tran.CoinTypeID == coin.ID {
					found = true
					break
				}
			}

			if !found {
				continue
			}

			state := txmgrpb.TxState_StateFail
			switch tran.State {
			case "created":
				state = txmgrpb.TxState_StateCreated
			case "wait":
				state = txmgrpb.TxState_StateWait
			case "paying":
				state = txmgrpb.TxState_StateTransferring
			case "successful":
				state = txmgrpb.TxState_StateSuccessful
			}

			txType := txmgrpb.TxType_TxWithdraw
			switch tran.CreatedFor {
			case "collecting":
				txType = txmgrpb.TxType_TxPaymentCollect
			case "withdraw":
			case "platform-benefit":
				txType = txmgrpb.TxType_TxBenefit
			case "compatible":
				if strings.Contains(tran.Message, "transfer gas") {
					txType = txmgrpb.TxType_TxFeedGas
				}
			}

			_, err := tx.
				Tran.
				Create().
				SetID(tran.ID).
				SetCoinTypeID(tran.CoinTypeID).
				SetFromAccountID(tran.FromAddressID).
				SetToAccountID(tran.ToAddressID).
				SetAmount(decimal.NewFromInt(int64(tran.Amount)).Div(decimal.NewFromInt(1000000000000))).
				SetFeeAmount(decimal.NewFromInt(int64(tran.TransactionFee)).Div(decimal.NewFromInt(1000000000000))).
				SetExtra(tran.Message).
				SetState(state.String()).
				SetType(txType.String()).
				Save(_ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func migrateBilling(ctx context.Context) error {
	conn, err := open(billingconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("migrateBilling", "error", err)
		return err
	}
	defer conn.Close()

	if err := migrateTx(ctx, conn); err != nil {
		logger.Sugar().Errorw("migrateBilling", "error", err)
		return err
	}

	return nil
}

var coinInfos = []*coininfoent.CoinInfo{}
var apps = []*appmwpb.App{}
var withdrawSettings = []*billingent.AppWithdrawSetting{}

func _migrateCoinInfo(ctx context.Context, conn *sql.DB) error { //nolint
	cli1 := coininfoent.NewClient(coininfoent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	coins, err := cli1.
		CoinInfo.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
		return err
	}

	coinInfos = coins

	conn1, err := open(projinfoconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
		return err
	}
	defer conn1.Close()

	cli2 := projinfoent.NewClient(projinfoent.Driver(entsql.OpenDB(dialect.MySQL, conn1)))
	prodInfos, err := cli2.
		CoinProductInfo.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
		return err
	}

	conn2, err := open(billingconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
		return err
	}
	defer conn2.Close()

	cli3 := billingent.NewClient(billingent.Driver(entsql.OpenDB(dialect.MySQL, conn2)))
	_withdrawSettings, err := cli3.
		AppWithdrawSetting.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
		return err
	}

	withdrawSettings = _withdrawSettings

	offset := int32(0)
	limit := int32(1000)

	for {
		_apps, _, err := appmwcli.GetApps(ctx, offset, limit)
		if err != nil {
			logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
			return err
		}
		if len(_apps) == 0 {
			break
		}

		apps = append(apps, _apps...)
		offset += limit
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		infos, err := tx.
			CoinBase.
			Query().
			Limit(1).
			All(ctx)
		if err != nil {
			logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, coin := range coins {
			_, err := tx.
				CoinBase.
				Create().
				SetID(coin.ID).
				SetName(coin.Name).
				SetUnit(coin.Unit).
				SetLogo(coin.Logo).
				SetPresale(coin.PreSale).
				SetEnv(coin.Env).
				SetReservedAmount(decimal.NewFromInt(int64(coin.ReservedAmount)).Div(decimal.NewFromInt(1000000000000))).
				SetForPay(coin.ForPay).
				Save(_ctx)
			if err != nil {
				return err
			}

			_, err = tx.
				CoinExtra.
				Create().
				SetCoinTypeID(coin.ID).
				SetHomePage(coin.HomePage).
				SetSpecs(coin.Specs).
				Save(_ctx)
			if err != nil {
				return err
			}

			for _, app := range apps {
				productPage := ""
				for _, prodInfo := range prodInfos {
					if app.ID == prodInfo.AppID.String() && coin.ID == prodInfo.CoinTypeID {
						productPage = prodInfo.ProductPage
						break
					}
				}

				withdrawAutoReviewAmount := decimal.NewFromInt(0)
				for _, setting := range withdrawSettings {
					if setting.AppID.String() == app.ID && coin.ID == setting.CoinTypeID {
						withdrawAutoReviewAmount = decimal.NewFromInt(int64(setting.WithdrawAutoReviewCoinAmount)).Div(decimal.NewFromInt(1000000000000))
					}
				}

				_, err := tx.
					AppCoin.
					Create().
					SetAppID(uuid.MustParse(app.ID)).
					SetCoinTypeID(coin.ID).
					SetName(coin.Name).
					SetLogo(coin.Logo).
					SetForPay(coin.ForPay).
					SetWithdrawAutoReviewAmount(withdrawAutoReviewAmount).
					SetProductPage(productPage).
					Save(_ctx)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func migrateCoinInfo(ctx context.Context) error {
	conn, err := open(coininfoconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinInfo", "error", err)
		return err
	}
	defer conn.Close()

	if err := _migrateCoinInfo(ctx, conn); err != nil {
		logger.Sugar().Errorw("migrateCoinInfo", "error", err)
		return err
	}

	return nil
}

func migrateCoinDescription(ctx context.Context, conn *sql.DB) error {
	cli1 := projinfoent.NewClient(projinfoent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	descs, err := cli1.
		CoinDescription.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinDescription", "error", err)
		return err
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		infos, err := tx.
			CoinDescription.
			Query().
			Limit(1).
			All(ctx)
		if err != nil {
			logger.Sugar().Errorw("migrateCoinDescription", "error", err)
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, desc := range descs {
			if desc.UsedFor != "PRODUCTDETAILS" {
				logger.Sugar().Warnw("migrateCoinDescription", "UsedFor", desc.UsedFor)
				continue
			}

			found := false
			for _, app := range apps {
				if app.ID == desc.AppID.String() {
					found = true
					break
				}
			}

			if !found {
				continue
			}

			found = false
			for _, coin := range coinInfos {
				if coin.ID == desc.CoinTypeID {
					found = true
					break
				}
			}

			if !found {
				continue
			}

			_, err := tx.
				CoinDescription.
				Create().
				SetAppID(desc.AppID).
				SetCoinTypeID(desc.CoinTypeID).
				SetUsedFor(descmgrpb.UsedFor_ProductPage.String()).
				SetTitle(desc.Title).
				SetMessage(desc.Message).
				Save(_ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func migrateProjectInfo(ctx context.Context) error {
	conn, err := open(projinfoconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("migrateProjectInfo", "error", err)
		return err
	}
	defer conn.Close()

	if err := migrateCoinDescription(ctx, conn); err != nil {
		logger.Sugar().Errorw("migrateProjectInfo", "error", err)
		return err
	}

	return nil
}

func migrateCurrency(ctx context.Context, conn *sql.DB) error {
	cli1 := oracleent.NewClient(oracleent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	currencies, err := cli1.
		Currency.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCurency", "error", err)
		return err
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		infos, err := tx.
			ExchangeRate.
			Query().
			Limit(1).
			All(ctx)
		if err != nil {
			logger.Sugar().Errorw("migrateCurrency", "error", err)
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, currency := range currencies {
			logger.Sugar().Infow("migrateCurrency", "Currency", currency)
			found := false
			for _, app := range apps {
				if app.ID == currency.AppID.String() {
					found = true
					break
				}
			}

			if !found {
				continue
			}

			found = false
			for _, coin := range coinInfos {
				if coin.ID == currency.CoinTypeID {
					found = true
					break
				}
			}

			if !found {
				continue
			}

			_, err := tx.
				ExchangeRate.
				Create().
				SetAppID(currency.AppID).
				SetCoinTypeID(currency.CoinTypeID).
				SetMarketValue(decimal.NewFromInt(int64(currency.PriceVsUsdt)).Div(decimal.NewFromInt(1000000000000))).
				SetSettleValue(decimal.NewFromInt(int64(currency.AppPriceVsUsdt)).Div(decimal.NewFromInt(1000000000000))).
				SetSettlePercent(90).
				Save(_ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func migrateOracle(ctx context.Context) error {
	conn, err := open(oracleconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("migrateOracle", "error", err)
		return err
	}
	defer conn.Close()

	if err := migrateCurrency(ctx, conn); err != nil {
		logger.Sugar().Errorw("migrateOracle", "error", err)
		return err
	}

	return nil
}

func migrateCoinGas(ctx context.Context, conn *sql.DB) error { //nolint
	cli1 := gasfeederent.NewClient(gasfeederent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	gases, err := cli1.
		CoinGas.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinGas", "error", err)
		return err
	}

	conn2, err := open(billingconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
		return err
	}
	defer conn2.Close()

	cli3 := billingent.NewClient(billingent.Driver(entsql.OpenDB(dialect.MySQL, conn2)))
	coinSettings, err := cli3.
		CoinSetting.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinGas", "error", err)
		return err
	}

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
		infos, err := tx.
			Setting.
			Query().
			Limit(1).
			All(ctx)
		if err != nil {
			logger.Sugar().Errorw("migrateCoinGas", "error", err)
			return err
		}
		if len(infos) > 0 {
			return nil
		}

		for _, coin := range coinInfos {
			feeCoinTypeID := coin.ID
			defaultFeeAmount := decimal.RequireFromString("0.001")
			collectAmount := decimal.RequireFromString("1000")

			logger.Sugar().Infow("migrateCoinGas", "Coin", coin.Name)

			for _, gas := range gases {
				if coin.ID == gas.CoinTypeID {
					feeCoinTypeID = gas.GasCoinTypeID
					defaultFeeAmount = decimal.NewFromInt(int64(gas.DepositAmount)).Div(decimal.NewFromInt(1000000000000))
					break
				}
			}

			for _, set := range coinSettings {
				if set.CoinTypeID == coin.ID {
					collectAmount = decimal.NewFromInt(int64(set.WarmAccountCoinAmount)).Div(decimal.NewFromInt(1000000000000))
					break
				}
			}

			_, err = tx.
				Setting.
				Create().
				SetCoinTypeID(coin.ID).
				SetFeeCoinTypeID(feeCoinTypeID).
				SetWithdrawFeeByStableUsd(true).
				SetWithdrawFeeAmount(decimal.NewFromInt(2)).
				SetCollectFeeAmount(defaultFeeAmount).
				SetHotWalletFeeAmount(defaultFeeAmount).
				SetLowFeeAmount(defaultFeeAmount).
				SetHotWalletAccountAmount(collectAmount).
				SetPaymentAccountCollectAmount(collectAmount).
				Save(_ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func migrateGasFeeder(ctx context.Context) error {
	conn, err := open(gasfeederconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("migrateGasFeeder", "error", err)
		return err
	}
	defer conn.Close()

	if err := migrateCoinGas(ctx, conn); err != nil {
		logger.Sugar().Errorw("migrateGasFeeder", "error", err)
		return err
	}

	return nil
}

func Migrate(ctx context.Context) error {
	if err := db.Init(); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	// Migrate coin info
	if err := migrateCoinInfo(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	// Migrate project info
	if err := migrateProjectInfo(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	// Migrate oracle
	if err := migrateOracle(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	// Migrate gas feeder
	if err := migrateGasFeeder(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	// Migrate billing
	if err := migrateBilling(ctx); err != nil {
		logger.Sugar().Errorw("Migrate", "error", err)
		return err
	}

	return nil
}
