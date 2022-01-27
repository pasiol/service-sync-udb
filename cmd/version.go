package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of service-sync-udb",
	Long:  `Print the version number of service-sync-udb`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service-sync-udb v0.1.0")
	},
}
