package main

import (
	"context"

	"github.com/NpoolPlatform/chain-gateway/pkg/migrator"
	"github.com/NpoolPlatform/go-service-framework/pkg/action"

	cli "github.com/urfave/cli/v2"
)

var migrateCmd = &cli.Command{
	Name:    "migrate",
	Aliases: []string{"m"},
	Usage:   "Migrate database",
	Action: func(c *cli.Context) error {
		err := action.Run(
			c.Context,
			func(ctx context.Context) error {
				return migrator.Migrate(ctx)
			},
			rpcRegister,
			rpcGatewayRegister,
			nil,
		)
		return err
	},
}

func runMigrate(ctx context.Context) error {
	if err := migrator.Migrate(ctx); err != nil {
		return err
	}
	return nil
}
