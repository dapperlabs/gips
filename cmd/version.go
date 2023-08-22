package cmd

import (
	"fmt"
	"log"

	"github.com/darron/gips/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information for " + binaryName,
	Run: func(cmd *cobra.Command, args []string) {
		version()
	},
}

func version() {
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	v := config.GetVersionInfo()
	conf.Logger.Info("Version Info", "info", fmt.Sprintf("%#v\n", v))
}
