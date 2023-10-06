package cmd

import (
	"context"

	"github.com/hosseintrz/gaterun/pkg/api"
	"github.com/hosseintrz/gaterun/pkg/gateway"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type GatewayOpts struct {
	ConfigFile    string
	RunMigrations bool
}

func NewStartGatewayCmd(ctx context.Context) *cobra.Command {
	opts := &GatewayOpts{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "start gaterun api gateway instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunGatewayStart(ctx, opts)
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.ConfigFile, "conf", "c", "gaterun.conf.yml", "default is './gaterun.conf.yml'")
	cmd.PersistentFlags().BoolVar(&opts.RunMigrations, "run-migrations", false, "run migrations before starting")

	return cmd
}

func RunGatewayStart(ctx context.Context, opts *GatewayOpts) error {
	initConfig(opts.ConfigFile)
	initLogger()
	initDatabase()

	log.Info("starting gaterun ")

	go func() {
		api.ServeApi()
	}()
	gw := gateway.NewGateway(globalConfig.ServiceConfig)
	gw.Start()

	return nil
}
