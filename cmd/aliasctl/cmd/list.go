package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	Long:  `List all aliases defined in the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		am.ListAliases()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
