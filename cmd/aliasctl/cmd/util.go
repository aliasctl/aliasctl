package cmd

import (
	"strings"
)

// parseAliasDefinition attempts to extract the alias name and command from a definition
func parseAliasDefinition(definition, shellType string) (name string, command string) {
	definition = strings.TrimSpace(definition)

	switch shellType {
	case "bash", "zsh", "fish", "ksh":
		if strings.HasPrefix(definition, "alias ") {
			parts := strings.SplitN(strings.TrimPrefix(definition, "alias "), "=", 2)
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				// Remove surrounding quotes
				command = strings.Trim(strings.TrimSpace(parts[1]), "'\"")
				return
			}
		}
	case "powershell", "pwsh":
		if strings.HasPrefix(definition, "Set-Alias ") {
			parts := strings.Fields(strings.TrimPrefix(definition, "Set-Alias "))
			if len(parts) >= 2 {
				name = parts[0]
				command = parts[1]
				return
			}
		} else if strings.HasPrefix(definition, "function ") {
			parts := strings.SplitN(strings.TrimPrefix(definition, "function "), " {", 2)
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				command = strings.TrimSpace(strings.TrimSuffix(parts[1], "}"))
				return
			}
		}
	case "cmd":
		if strings.HasPrefix(definition, "doskey ") {
			parts := strings.SplitN(strings.TrimPrefix(definition, "doskey "), "=", 2)
			if len(parts) == 2 {
				name = strings.TrimSpace(parts[0])
				command = strings.TrimSpace(parts[1])
				return
			}
		}
	}

	// Fallback: if we couldn't parse with the specific shell format,
	// try to use a generic approach - look for the first space or equals
	if strings.Contains(definition, "=") {
		parts := strings.SplitN(definition, "=", 2)
		if len(parts) == 2 {
			name = strings.TrimSpace(parts[0])
			command = strings.TrimSpace(parts[1])
			return
		}
	}

	return "", ""
}
