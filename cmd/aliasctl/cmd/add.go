package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// addCmd represents the add command which creates a new alias and saves it to storage.
// It takes a name and a command as arguments, joining multiple command arguments into a single string.
// Example usage: aliasctl add ll "ls -la"
var addCmd = &cobra.Command{
	Use:   "add [name] [command]",
	Short: "Add a new alias",
	Long:  `Add a new alias mapping a name to a shell command.`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		command := strings.Join(args[1:], " ")

		am.AddAlias(name, command)
		if err := am.SaveAliases(); err != nil {
			return fmt.Errorf("failed to save alias: %w\n\nTry ensuring you have write permissions to %s or specify an alternative location with 'aliasctl set-file'", err, am.AliasStore)
		}

		fmt.Printf("Added alias: %s = %s\n", name, command)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
