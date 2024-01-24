package cmd

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/hosseintrz/gaterun/config"
	"github.com/hosseintrz/gaterun/config/persistence"
	"github.com/hosseintrz/gaterun/pkg/api"
	"github.com/hosseintrz/gaterun/pkg/gateway"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

type GatewayOpts struct {
	ConfigFile    string
	RunMigrations bool
	ReadEndpoints bool
	DBMode        bool
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
	cmd.PersistentFlags().BoolVar(&opts.ReadEndpoints, "read-endpoints", false, "import endpoints from config file")
	cmd.PersistentFlags().BoolVarP(&opts.DBMode, "dbmode", "d", false, "if set to false, configs will be read from conf or else from db")

	return cmd
}

func RunGatewayStart(ctx context.Context, opts *GatewayOpts) error {
	initConfig(opts.ConfigFile)
	initLogger()
	initDatabase()
	initRedis(ctx)

	persisted := false
	if opts.DBMode {
		conf, err := persistence.AssembleConfig(ctx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Info("no configs found in db")
			} else {
				return err
			}
		} else {
			persisted = true
			globalConfig.ServiceConfig = *conf
		}
	}

	config.SetGlobalConf(globalConfig)

	if !persisted {
		err := persistence.SaveConfigToDB(ctx, globalConfig)
		if err != nil {
			return err
		}
	}

	go func() {
		api.ServeApi()
	}()

	gw := gateway.NewGateway(globalConfig.ServiceConfig)
	go gw.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	return nil
}
