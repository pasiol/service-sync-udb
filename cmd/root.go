package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wilma-users",
	Short: "wilma-user is command line tool manage Wilma accounts",
	Long: `wilma-user is command line tool manage Wilma accounts.
               Syncing personal ids from Primus to user database`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
