package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Sehsyha/crounch-back/configuration"
)

const (
	parameterMock = "mock"
	defaultMock   = false

	// Database
	parameterDBConnectionURI = "db-connection-uri"
	defaultDBConnectionURI   = "postgresql://postgres:password@database/postgres?sslmode=disable"
)

var (
	config = &configuration.Config{}

	rootCmd = &cobra.Command{
		Use:   "crounch",
		Short: "crounch",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func initConfig() {
	viper.AutomaticEnv()

	config.Mock = viper.GetBool(parameterMock)

	config.DBConnectionURI = viper.GetString(parameterDBConnectionURI)
}
