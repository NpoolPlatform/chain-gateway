//nolint:nolintlint
package migrator

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"

	servicename "github.com/NpoolPlatform/chain-gateway/pkg/servicename"
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

	return nil
}
