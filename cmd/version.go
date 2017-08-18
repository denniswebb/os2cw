package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	//BuildVersion is the current version of the binary.  Set at build time.
	BuildVersion = "0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
}

func init() {
	cmd := versionCmd
	rootCommand.AddCommand(cmd)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s Version %s\n", appName, BuildVersion)
	}
}
