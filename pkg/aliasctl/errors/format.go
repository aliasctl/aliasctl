package errors

import (
	"fmt"
	"strings"
)

// Format creates a well-formatted error message with optional hints
func Format(msg string, err error, hints ...string) error {
	var sb strings.Builder

	// Main error message
	sb.WriteString(msg)

	if err != nil {
		// Original error
		sb.WriteString(": ")
		sb.WriteString(err.Error())
	}

	// Add hints if provided
	if len(hints) > 0 {
		sb.WriteString("\n\n")
		for _, hint := range hints {
			sb.WriteString(hint)
			sb.WriteString("\n")
		}
	}

	return fmt.Errorf(sb.String())
}

// FormatNetworkError formats a network-related error with appropriate hints
func FormatNetworkError(msg string, err error) error {
	return Format(msg, err,
		"Possible causes:",
		"- The service might not be running",
		"- Network connectivity issues",
		"- Incorrect endpoint URL")
}

// FormatPermissionError formats a permission-related error with appropriate hints
func FormatPermissionError(path string, err error) error {
	return Format(fmt.Sprintf("Permission denied for %s", path), err,
		"Possible solutions:",
		"- Check if you have appropriate file/directory permissions",
		"- Try running with elevated privileges",
		"- Specify an alternative location with 'aliasctl set-file'")
}

// FormatConfigError formats a configuration-related error with appropriate hints
func FormatConfigError(msg string, err error) error {
	return Format(msg, err,
		"Possible solutions:",
		"- Check your configuration file format",
		"- Consider resetting configuration with `set-file` or `set-shell`",
		"- Ensure the configuration directory exists and is writable")
}

// FormatNotFoundError formats a not found error with appropriate hints
func FormatNotFoundError(resourceType string, name string, suggestion string) error {
	msg := fmt.Sprintf("%s '%s' not found", resourceType, name)
	hints := []string{
		fmt.Sprintf("Suggestion: %s", suggestion),
	}
	return Format(msg, nil, hints...)
}

// FormatAPIError formats an API-related error with appropriate hints
func FormatAPIError(provider string, err error) error {
	hints := []string{
		"Possible causes:",
	}

	// Add provider-specific hints
	switch provider {
	case "ollama":
		hints = append(hints,
			"- Ollama service might not be running (start with 'ollama serve')",
			"- The specified model might not be downloaded (try 'ollama pull <model>')",
			"- Incorrect Ollama endpoint URL")
	case "openai":
		hints = append(hints,
			"- API key might be invalid or expired",
			"- The model name might be incorrect",
			"- You may have reached your API usage limit",
			"- Incorrect OpenAI endpoint URL")
	case "anthropic":
		hints = append(hints,
			"- API key might be invalid or expired",
			"- The model name might be incorrect",
			"- You may have reached your API usage limit",
			"- Incorrect Anthropic endpoint URL")
	default:
		hints = append(hints,
			"- API key might be invalid",
			"- Service might be unavailable",
			"- Network connectivity issues")
	}

	return Format(fmt.Sprintf("%s API error", provider), err, hints...)
}
