package main

import (
	"fmt"
	"os"

	servicename "github.com/NpoolPlatform/chain-gateway/pkg/servicename"

	"github.com/NpoolPlatform/go-service-framework/pkg/app"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	mysqlconst "github.com/NpoolPlatform/go-service-framework/pkg/mysql/const"
	rabbitmqconst "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/const"
	redisconst "github.com/NpoolPlatform/go-service-framework/pkg/redis/const"

	billingconst "github.com/NpoolPlatform/cloud-hashing-billing/pkg/message/const"
	gasfeederconst "github.com/NpoolPlatform/gas-feeder/pkg/message/const"
	oracleconst "github.com/NpoolPlatform/oracle-manager/pkg/message/const"
	projinfoconst "github.com/NpoolPlatform/project-info-manager/pkg/message/const"
	coininfoconst "github.com/NpoolPlatform/sphinx-coininfo/pkg/message/const"

	cli "github.com/urfave/cli/v2"
)

func main() {
	commands := cli.Commands{
		runCmd,
	}

	description := fmt.Sprintf("my %v service cli\nFor help on any individual command run <%v COMMAND -h>\n",
		servicename.ServiceName, servicename.ServiceName)
	err := app.Init(
		servicename.ServiceName,
		description,
		"",
		"",
		"./",
		nil,
		commands,
		mysqlconst.MysqlServiceName,
		rabbitmqconst.RabbitMQServiceName,
		redisconst.RedisServiceName,
		billingconst.ServiceName,
		coininfoconst.ServiceName,
		projinfoconst.ServiceName,
		oracleconst.ServiceName,
		gasfeederconst.ServiceName,
	)
	if err != nil {
		logger.Sugar().Errorf("fail to create %v: %v", servicename.ServiceName, err)
		return
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Sugar().Errorf("fail to run %v: %v", servicename.ServiceName, err)
	}
}
