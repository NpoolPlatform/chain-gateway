package main

import (
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
			migrator.Migrate,
			rpcRegister,
			rpcGatewayRegister,
			nil,
		)
		return err
	},
}
