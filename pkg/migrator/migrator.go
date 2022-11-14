package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NpoolPlatform/chain-manager/pkg/db"
	"github.com/NpoolPlatform/chain-manager/pkg/db/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/shopspring/decimal"

	billingent "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent"
	entcoinsetting "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent/coinsetting"

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

func migrateCoinSetting(ctx context.Context, conn *sql.DB) error {
	cli1 := billingent.NewClient(billingent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	settings, err := cli1.
		CoinSetting.
		Query().
		Where(
			entcoinsetting.DeleteAt(0),
		).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinSetting", "error", err)
		return err
	}

	for _, setting := range settings {
		logger.Sugar().Infow("migrateCoinSetting", "Setting", setting)
	}

	return nil
}

func migrateBilling(ctx context.Context) error {
	conn, err := open(billingconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("migrateBilling", "error", err)
		return err
	}
	defer conn.Close()

	if err := migrateCoinSetting(ctx, conn); err != nil {
		logger.Sugar().Errorw("migrateBilling", "error", err)
		return err
	}

	return nil
}

var coinInfos = []*coininfoent.CoinInfo{}
var coinProductInfos = []*projinfoent.CoinProductInfo{}
var coinDescs = []*projinfoent.CoinDescription{}
var apps = []*appmwpb.App{}
var withdrawSetings = []*billingent.AppWithdrawSetting{}

func _migrateCoinInfo(ctx context.Context, conn *sql.DB) error {
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
	withdrawSettings, err := cli3.
		AppWithdrawSetting.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
		return err
	}

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

	cli, err := db.Client()
	if err != nil {
		logger.Sugar().Errorw("_migrateCoinInfo", "error", err)
		return err
	}
	infos, err := cli.
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

	return db.WithTx(ctx, func(_ctx context.Context, tx *ent.Tx) error {
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

func migrateCoinProductInfo(ctx context.Context, conn *sql.DB) error {
	cli1 := projinfoent.NewClient(projinfoent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	prodInfos, err := cli1.
		CoinProductInfo.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinProductInfo", "error", err)
		return err
	}

	coinProductInfos = prodInfos

	cli, err := db.Client()
	if err != nil {
		logger.Sugar().Errorw("migrateCoinProductInfo", "error", err)
		return err
	}
	infos, err := cli.
		AppCoin.
		Query().
		Limit(1).
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinProductInfo", "error", err)
		return err
	}
	if len(infos) > 0 {
		return nil
	}

	for _, info := range infos {
		logger.Sugar().Infow("migrateCoinProductInfo", "CoinProductInfo", info)
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

	for _, desc := range descs {
		logger.Sugar().Infow("migrateCoinDescription", "CoinDescription", desc)
	}

	return nil
}

func migrateProjectInfo(ctx context.Context) error {
	conn, err := open(projinfoconst.ServiceName)
	if err != nil {
		logger.Sugar().Errorw("migrateProjectInfo", "error", err)
		return err
	}
	defer conn.Close()

	if err := migrateCoinProductInfo(ctx, conn); err != nil {
		logger.Sugar().Errorw("migrateProjectInfo", "error", err)
		return err
	}

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

	for _, currency := range currencies {
		logger.Sugar().Infow("migrateCurrency", "Currency", currency)
	}

	return nil
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

func migrateCoinGas(ctx context.Context, conn *sql.DB) error {
	cli1 := gasfeederent.NewClient(gasfeederent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	gases, err := cli1.
		CoinGas.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinGas", "error", err)
		return err
	}

	for _, gas := range gases {
		logger.Sugar().Infow("migrateCoinGas", "Gas", gas)
	}

	return nil
}

func migrateDeposit(ctx context.Context, conn *sql.DB) error {
	cli1 := gasfeederent.NewClient(gasfeederent.Driver(entsql.OpenDB(dialect.MySQL, conn)))
	deposits, err := cli1.
		Deposit.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateDeposit", "error", err)
		return err
	}

	for _, deposit := range deposits {
		logger.Sugar().Infow("migrateDeposit", "Deposit", deposit)
	}

	return nil
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

	if err := migrateDeposit(ctx, conn); err != nil {
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

	// Migrate billing
	if err := migrateBilling(ctx); err != nil {
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

	return nil
}
