package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `Generate completion script for the specified shell.
The script can be integrated with your shell to enable tab completion for aliasctl commands.`,
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		shellType := args[0]

		var err error
		switch shellType {
		case "bash":
			err = rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			err = rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			err = rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}

		if err != nil {
			return fmt.Errorf("failed to generate %s completion script: %w\n\nTry running with administrator/root privileges if needed", shellType, err)
		}

		return nil
	},
}

var installCompletionCmd = &cobra.Command{
	Use:   "install-completion",
	Short: "Install completion script for current shell",
	Long:  `Install shell completion script for the current shell.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := am.InstallCompletionScript(); err != nil {
			// Get the target path for better error message
			homeDir, _ := os.UserHomeDir()
			var targetPath string

			switch am.Shell {
			case "bash":
				targetPath = filepath.Join(homeDir, ".bash_completion.d")
			case "zsh":
				targetPath = filepath.Join(homeDir, ".zsh", "completion")
			case "fish":
				targetPath = filepath.Join(homeDir, ".config", "fish", "completions")
			case "powershell", "pwsh":
				targetPath = "PowerShell profile directory"
			}

			return fmt.Errorf("failed to install completion script: %w\n\nEnsure you have write permissions to %s or run with administrator/root privileges", err, targetPath)
		}
		fmt.Println("Shell completion script successfully installed")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	rootCmd.AddCommand(installCompletionCmd)
}
