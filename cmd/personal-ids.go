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
	Short: "Syncing missing personal ids to SQL database.",
	Long: `Syncing missing personal ids to SQL database.
	Retrieves updated information from Primus.
	Updates to SQL-database personal id	and 
	student id-fields, using primus id as foreign key.
	
	syncudb personal-ids`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.SyncIds()
	},
}
