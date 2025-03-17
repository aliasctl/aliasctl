package aliasctl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ApplyAliases writes the aliases to the shell configuration file.
func (am *AliasManager) ApplyAliases() error {
	existingContent := ""
	existingAliasSection := false

	if _, err := os.Stat(am.AliasFile); err == nil {
		content, err := os.ReadFile(am.AliasFile)
		if err == nil {
			existingContent = string(content)
			if strings.Contains(existingContent, "# Aliases managed by AliasCtl") {
				existingAliasSection = true
			}
		}
	}

	var newContent strings.Builder

	if !existingAliasSection && existingContent != "" {
		newContent.WriteString(existingContent)
		if !strings.HasSuffix(existingContent, "\n") {
			newContent.WriteString("\n")
		}
		newContent.WriteString("\n# Aliases managed by AliasCtl\n")
	} else if existingAliasSection {
		parts := strings.SplitN(existingContent, "# Aliases managed by AliasCtl", 2)
		newContent.WriteString(parts[0])
		newContent.WriteString("# Aliases managed by AliasCtl\n")
		if len(parts) > 1 && strings.Contains(parts[1], "# End of aliases managed by AliasCtl") {
			// Do nothing - we'll add a new section
		}
	} else {
		newContent.WriteString("# Aliases managed by AliasCtl\n")
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
			switch am.Shell {
			case ShellPowerShell, ShellPowerShellCore:
				if strings.Contains(command, " ") {
					newContent.WriteString(fmt.Sprintf("function %s { %s }\n", name, command))
				} else {
					newContent.WriteString(fmt.Sprintf("Set-Alias %s %s\n", name, command))
				}
			case ShellCmd:
				newContent.WriteString(fmt.Sprintf("doskey %s=%s\n", name, command))
			case ShellFish:
				if strings.Contains(command, " ") {
					newContent.WriteString(fmt.Sprintf("function %s\n    %s\nend\n", name, command))
				} else {
					newContent.WriteString(fmt.Sprintf("alias %s '%s'\n", name, command))
				}
			default:
				newContent.WriteString(fmt.Sprintf("alias %s='%s'\n", name, command))
			}
		}
	}

	newContent.WriteString("# End of aliases managed by AliasCtl\n")

	if existingAliasSection && strings.Contains(existingContent, "# End of aliases managed by AliasCtl") {
		parts := strings.SplitN(existingContent, "# End of aliases managed by AliasCtl", 2)
		if len(parts) > 1 {
			newContent.WriteString(parts[1])
		}
	}

	dir := filepath.Dir(am.AliasFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(am.AliasFile, []byte(newContent.String()), 0644)
}

// ImportAliasesFromShell imports aliases from the shell configuration file.
func (am *AliasManager) ImportAliasesFromShell() error {
	file, err := os.Open(am.AliasFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		switch am.Shell {
		case ShellPowerShell, ShellPowerShellCore:
			if strings.HasPrefix(line, "function ") {
				parts := strings.SplitN(line[9:], " ", 2)
				if len(parts) == 2 && strings.Contains(parts[1], "{") {
					name := parts[0]
					cmdParts := strings.SplitN(parts[1], "{", 2)
					if len(cmdParts) == 2 {
						command := strings.TrimSpace(cmdParts[1])
						command = strings.TrimSuffix(command, "}")
						commands := am.Aliases[name]
						switch am.Shell {
						case ShellPowerShell:
							commands.PowerShell = strings.TrimSpace(command)
						case ShellPowerShellCore:
							commands.PowerShellCore = strings.TrimSpace(command)
						}
						am.Aliases[name] = commands
					}
				}
			} else if strings.HasPrefix(line, "Set-Alias ") {
				parts := strings.Fields(line[10:])
				if len(parts) >= 2 {
					commands := am.Aliases[parts[0]]
					switch am.Shell {
					case ShellPowerShell:
						commands.PowerShell = parts[1]
					case ShellPowerShellCore:
						commands.PowerShellCore = parts[1]
					}
					am.Aliases[parts[0]] = commands
				}
			}
		case ShellCmd:
			if strings.HasPrefix(line, "doskey ") {
				parts := strings.SplitN(line[7:], "=", 2)
				if len(parts) == 2 {
					commands := am.Aliases[parts[0]]
					commands.Cmd = parts[1]
					am.Aliases[parts[0]] = commands
				}
			}
		case ShellFish:
			if strings.HasPrefix(line, "alias ") {
				line = strings.TrimPrefix(line, "alias ")
				parts := strings.SplitN(line, " ", 2)
				if len(parts) == 2 {
					name := parts[0]
					command := strings.Trim(parts[1], "'\"")
					commands := am.Aliases[name]
					commands.Fish = command
					am.Aliases[name] = commands
				}
			} else if strings.HasPrefix(line, "function ") {
				parts := strings.SplitN(line[9:], " ", 2)
				if len(parts) >= 1 {
					name := strings.TrimSuffix(parts[0], ";")
					if scanner.Scan() {
						command := strings.TrimSpace(scanner.Text())
						if !strings.HasPrefix(command, "end") {
							commands := am.Aliases[name]
							commands.Fish = command
							am.Aliases[name] = commands
						}
					}
				}
			}
		default:
			if strings.HasPrefix(line, "alias ") {
				line = strings.TrimPrefix(line, "alias ")
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					name := parts[0]
					command := strings.Trim(parts[1], "'\"")
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
					}
					am.Aliases[name] = commands
				}
			}
		}
	}

	return am.SaveAliases()
}

// ExportAliases exports aliases to a different shell format.
func (am *AliasManager) ExportAliases(targetShell, outputFile string) error {
	var content strings.Builder
	content.WriteString("# Aliases exported by AliasCtl\n")

	for name, commands := range am.Aliases {
		var command string
		switch targetShell {
		case "bash":
			command = commands.Bash
		case "zsh":
			command = commands.Zsh
		case "fish":
			command = commands.Fish
		case "ksh":
			command = commands.Ksh
		case "powershell":
			command = commands.PowerShell
		case "pwsh":
			command = commands.PowerShellCore
		case "cmd":
			command = commands.Cmd
		}

		if command != "" {
			if am.AIConfigured && string(am.Shell) != targetShell {
				convertedCommand, err := am.ConvertAlias(name, targetShell, "")
				if err == nil {
					command = convertedCommand
				}
			}

			switch targetShell {
			case "powershell", "pwsh":
				if strings.Contains(command, " ") {
					content.WriteString(fmt.Sprintf("function %s { %s }\n", name, command))
				} else {
					content.WriteString(fmt.Sprintf("Set-Alias %s %s\n", name, command))
				}
			case "cmd":
				content.WriteString(fmt.Sprintf("doskey %s=%s\n", name, command))
			case "fish":
				if strings.Contains(command, " ") {
					content.WriteString(fmt.Sprintf("function %s\n    %s\nend\n", name, command))
				} else {
					content.WriteString(fmt.Sprintf("alias %s '%s'\n", name, command))
				}
			default:
				content.WriteString(fmt.Sprintf("alias %s='%s'\n", name, command))
			}
		}
	}

	content.WriteString("# End of exported aliases\n")

	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(outputFile, []byte(content.String()), 0644)
}
