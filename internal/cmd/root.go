package cmd

import (
	"bank-authentication-system/internal/cmd/api"
	"bank-authentication-system/internal/cmd/consumer"
	"bank-authentication-system/internal/cmd/migrate"
	"bank-authentication-system/pkg/config"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

const exitFailure = 1

func Execute() {
	cfg := config.InitConfig()

	var cmd = &cobra.Command{
		Use:   "bank-authentication-system",
		Short: "bank",
	}

	logrus.Debugf("config loaded: %+v", cfg)

	migrate.Register(cmd, cfg)
	api.Register(cmd, cfg)
	consumer.Register(cmd, cfg)

	if err := cmd.Execute(); err != nil {
		logrus.Error(err.Error())
		os.Exit(exitFailure)
	}
}
