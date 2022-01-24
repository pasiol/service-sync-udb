package cmd

import (
	"github.com/pasiol/service-sync-udb/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(personalIdsCmd)
}

var personalIdsCmd = &cobra.Command{
	Use:   "personal-ids",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		internal.SyncIds()
	},
}
