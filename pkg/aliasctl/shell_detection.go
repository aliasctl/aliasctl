package aliasctl

import (
	"os"
	"path/filepath"
	"strings"
)

// DetectShellAndAliasFile determines the current shell and appropriate alias file.
func DetectShellAndAliasFile(platform string) (ShellType, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", ""
	}

	shellEnv := os.Getenv("SHELL")

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

	return shell, aliasFile
}
