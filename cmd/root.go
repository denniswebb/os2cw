package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	log "github.com/Sirupsen/logrus"

)

type RootCmd struct {
	cobraCommand *cobra.Command
}

var rootCommand = RootCmd{
	cobraCommand: &cobra.Command{
		Use: "os2cw",
		Short: "os2cw pushes select OS metrics to CloudWatch",
	},
}

func init() {
	cmd := rootCommand.cobraCommand

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}

func Execute() {
	flag.CommandLine.Parse([]string{})
	if err := rootCommand.cobraCommand.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}

