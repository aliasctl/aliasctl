package aliasctl

import (
	"strings"
)

// ParseAliasCommand attempts to extract an alias command from AI response text
func ParseAliasCommand(responseText string) string {
	// Clean up the response
	responseText = strings.TrimSpace(responseText)

	// Extract just the command if possible
	lines := strings.Split(responseText, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "alias ") ||
			strings.HasPrefix(line, "function ") ||
			strings.HasPrefix(line, "Set-Alias ") ||
			strings.HasPrefix(line, "doskey ") {
			return line
		}
	}

	return responseText
}
