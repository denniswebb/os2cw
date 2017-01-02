package cmd

import (
	"flag"
	"fmt"

	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	appName = "os2cw"
)

var (
	configFile  string
	rootCommand = &cobra.Command{
		Use:   appName,
		Short: fmt.Sprintf("%s pushes select OS metrics to CloudWatch", appName),
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCommand.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is %s.yaml)", appName))
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	viper.SetConfigName(appName)
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %", viper.ConfigFileUsed())
	}
}

func Execute() {
	flag.CommandLine.Parse([]string{})
	if err := rootCommand.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
