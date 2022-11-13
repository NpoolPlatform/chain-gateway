package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NpoolPlatform/chain-manager/pkg/db"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	billingent "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent"
	entcoinsetting "github.com/NpoolPlatform/cloud-hashing-billing/pkg/db/ent/coinsetting"

	projinfoent "github.com/NpoolPlatform/project-info-manager/pkg/db/ent"
	coininfoent "github.com/NpoolPlatform/sphinx-coininfo/pkg/db/ent"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	projinfoconst "github.com/NpoolPlatform/project-info-manager/pkg/message/const"
	coininfoconst "github.com/NpoolPlatform/sphinx-coininfo/pkg/message/const"

	_ "github.com/NpoolPlatform/project-info-manager/pkg/db/ent/runtime"
	_ "github.com/NpoolPlatform/sphinx-coininfo/pkg/db/ent/runtime"
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

	for _, coin := range coins {
		logger.Sugar().Infow("_migrateCoinInfo", "Coin", coin)
	}

	return nil
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
	infos, err := cli1.
		CoinProductInfo.
		Query().
		All(ctx)
	if err != nil {
		logger.Sugar().Errorw("migrateCoinProductInfo", "error", err)
		return err
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
	// Migrate gas feeder
	return nil
}
