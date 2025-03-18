package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command which deletes an existing alias.
// It requires exactly one argument - the name of the alias to remove.
// The command will return an error if the alias doesn't exist.
// Example usage: aliasctl remove ll
var removeCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove an alias",
	Long:  `Remove an existing alias by name.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if am.RemoveAlias(name) {
			if err := am.SaveAliases(); err != nil {
				return fmt.Errorf("alias '%s' was removed from memory but could not be saved to disk: %w\n\nTry checking if you have write permissions to %s", name, err, am.AliasStore)
			}
			fmt.Printf("Removed alias: %s\n", name)
			return nil
		}

		return fmt.Errorf("alias '%s' not found. Run 'aliasctl list' to see all available aliases", name)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
