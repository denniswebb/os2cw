package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	log "github.com/Sirupsen/logrus"
)

var (
	// value overwritten during build. This can be used to resolve issues.
	BuildVersion = "0.1"
)

type VersionCmd struct {
	cobraCommand *cobra.Command
}

var versionCmd = VersionCmd{
	cobraCommand: &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
	},
}

func init() {
	cmd := versionCmd.cobraCommand
	rootCommand.cobraCommand.AddCommand(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := versionCmd.Run()
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func (c *VersionCmd) Run() error {
	fmt.Printf("Version %s\n", BuildVersion)
	return nil
}
