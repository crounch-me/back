package cmd

import (
	"github.com/crounch-me/back/router"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:    "serve",
	Short:  "Serve endpoints",
	PreRun: initServeBindingFlags,
	Run: func(cmd *cobra.Command, args []string) {
		initConfig()

		log.
			WithField(parameterMock, config.Mock).
			WithField(parameterDBConnectionURI, config.DBConnectionURI).
			Warn("Configuration")

		router.Start(config)
	},
}

func init() {
	serveCmd.Flags().Bool(parameterMock, defaultMock, "Use this flag to mock external services (such as db)")
	serveCmd.Flags().String(parameterDBConnectionURI, defaultDBConnectionURI, "Use this flag to set the postgresql connection URI")
	serveCmd.Flags().String(parameterDBSchema, defaultDBSchema, "Use this flag to set database schema")
}

func initServeBindingFlags(cmd *cobra.Command, args []string) {
	_ = viper.BindPFlag(parameterMock, cmd.Flags().Lookup(parameterMock))
	_ = viper.BindPFlag(parameterDBConnectionURI, cmd.Flags().Lookup(parameterDBConnectionURI))
	_ = viper.BindPFlag(parameterDBSchema, cmd.Flags().Lookup(parameterDBSchema))
}
