package aliasctl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LoadAliases loads aliases from the alias store file.
func (am *AliasManager) LoadAliases() error {
	data, err := os.ReadFile(am.AliasStore)
	if err != nil {
		if os.IsNotExist(err) {
			am.Aliases = make(map[string]AliasCommands)
			return nil
		}
		return err
	}

	// Use the TOML support function if it exists, otherwise fall back to JSON
	if err := am.AddLoadAliasesTomlSupport(data); err != nil {
		// Fall back to JSON parsing for backward compatibility or if TOML parsing fails
		return json.Unmarshal(data, &am.Aliases)
	}

	return nil
}

// SaveAliases saves aliases to the alias store file.
func (am *AliasManager) SaveAliases() error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(am.AliasStore)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Try to use TOML support if available
	if err := am.AddSaveAliasesTomlSupport(); err == nil {
		return nil
	}

	// Fall back to JSON if TOML encoding fails
	data, err := json.MarshalIndent(am.Aliases, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(am.AliasStore, data, 0644)
}

// AddAlias adds a new alias.
func (am *AliasManager) AddAlias(name, command string) {
	commands := am.Aliases[name]
	switch am.Shell {
	case ShellBash:
		commands.Bash = command
	case ShellZsh:
		commands.Zsh = command
	case ShellFish:
		commands.Fish = command
	case ShellKsh:
		commands.Ksh = command
	case ShellPowerShell:
		commands.PowerShell = command
	case ShellPowerShellCore:
		commands.PowerShellCore = command
	case ShellCmd:
		commands.Cmd = command
	}
	am.Aliases[name] = commands
}

// RemoveAlias removes an alias.
func (am *AliasManager) RemoveAlias(name string) bool {
	if _, exists := am.Aliases[name]; exists {
		delete(am.Aliases, name)
		return true
	}
	return false
}

// ListAliases prints all aliases.
func (am *AliasManager) ListAliases() {
	fmt.Printf("Aliases for %s shell on %s platform:\n", am.Shell, am.Platform)
	if len(am.Aliases) == 0 {
		fmt.Println("No aliases defined.")
		return
	}

	for name, commands := range am.Aliases {
		var command string
		switch am.Shell {
		case ShellBash:
			command = commands.Bash
		case ShellZsh:
			command = commands.Zsh
		case ShellFish:
			command = commands.Fish
		case ShellKsh:
			command = commands.Ksh
		case ShellPowerShell:
			command = commands.PowerShell
		case ShellPowerShellCore:
			command = commands.PowerShellCore
		case ShellCmd:
			command = commands.Cmd
		}
		if command != "" {
			fmt.Printf("%s = %s\n", name, command)
		}
	}
}

// SetShell manually sets the shell type.
func (am *AliasManager) SetShell(shell string) error {
	switch shell {
	case "bash":
		am.Shell = ShellBash
	case "zsh":
		am.Shell = ShellZsh
	case "fish":
		am.Shell = ShellFish
	case "ksh":
		am.Shell = ShellKsh
	case "powershell":
		am.Shell = ShellPowerShell
	case "pwsh":
		am.Shell = ShellPowerShellCore
	case "cmd":
		am.Shell = ShellCmd
	default:
		return fmt.Errorf("unsupported shell: %s", shell)
	}
	return am.SaveConfig()
}

// SetAliasFile manually sets the alias file path.
func (am *AliasManager) SetAliasFile(filePath string) error {
	am.AliasFile = filePath
	return am.SaveConfig()
}
