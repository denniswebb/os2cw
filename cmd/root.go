package cmd

import (
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
	rootCommand.PersistentFlags().StringP("region", "", "", "AWS region")
	rootCommand.PersistentFlags().StringP("access-key", "", "", "AWS access key id")
	rootCommand.PersistentFlags().StringP("secret-key", "", "", "AWS secret access key")

	viper.BindPFlag("region", rootCommand.PersistentFlags().Lookup("region"))
	viper.BindPFlag("accessKey", rootCommand.PersistentFlags().Lookup("access-key"))
	viper.BindPFlag("secretKey", rootCommand.PersistentFlags().Lookup("secret-key"))

	viper.BindEnv("region", "AWS_REGION")
	viper.BindEnv("accessKey", "AWS_ACCESS_KEY_ID")
	viper.BindEnv("secretKey", "AWS_SECRET_ACCESS_KEY")
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
	//flag.CommandLine.Parse([]string{})
	if err := rootCommand.Execute(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
