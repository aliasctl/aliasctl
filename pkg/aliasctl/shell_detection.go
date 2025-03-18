package aliasctl

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DetectShellAndAliasFile determines the current shell and appropriate alias file.
func DetectShellAndAliasFile(platform string) (ShellType, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Warning: Failed to determine home directory: %v\n", err)
		fmt.Println("Using current directory as a fallback.")
		homeDir, _ = os.Getwd()
	}

	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" && platform != "windows" {
		fmt.Println("Warning: SHELL environment variable not set. Defaulting to bash.")
	}

	var shell ShellType
	var aliasFile string

	switch platform {
	case "windows":
		_, pwshErr := os.Stat(filepath.Join(os.Getenv("ProgramFiles"), "PowerShell", "7"))
		if pwshErr == nil {
			shell = ShellPowerShellCore
			aliasFile = filepath.Join(homeDir, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
		} else {
			shell = ShellPowerShell
			aliasFile = filepath.Join(homeDir, "Documents", "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1")
		}
	default:
		switch {
		case strings.Contains(shellEnv, "zsh"):
			shell = ShellZsh
			aliasFile = filepath.Join(homeDir, ".zshrc")
		case strings.Contains(shellEnv, "fish"):
			shell = ShellFish
			aliasFile = filepath.Join(homeDir, ".config", "fish", "config.fish")
		case strings.Contains(shellEnv, "ksh"):
			shell = ShellKsh
			aliasFile = filepath.Join(homeDir, ".kshrc")
		default:
			shell = ShellBash
			aliasFile = filepath.Join(homeDir, ".bash_aliases")
		}
	}

	// Verify the file exists or is writable
	if _, err := os.Stat(aliasFile); os.IsNotExist(err) {
		// Check if directory exists
		dir := filepath.Dir(aliasFile)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("Note: Destination directory %s does not exist. It will be created when needed.\n", dir)
		}
	}

	return shell, aliasFile
}
