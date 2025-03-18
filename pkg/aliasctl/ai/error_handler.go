package ai

import (
	"fmt"
	"strings"
)

// FormatAIError provides consistent formatting for AI-related errors
func FormatAIError(provider, errorType string, err error, details ...string) error {
	var message strings.Builder

	switch errorType {
	case "connection":
		message.WriteString(fmt.Sprintf("failed to connect to %s service", provider))
	case "authentication":
		message.WriteString(fmt.Sprintf("%s API authentication error", provider))
	case "model":
		message.WriteString(fmt.Sprintf("%s model not found", provider))
	case "response":
		message.WriteString(fmt.Sprintf("invalid response from %s API", provider))
	case "request":
		message.WriteString(fmt.Sprintf("failed to create %s request", provider))
	default:
		message.WriteString(fmt.Sprintf("%s error", provider))
	}

	if err != nil {
		message.WriteString(": ")
		message.WriteString(err.Error())
	}

	if len(details) > 0 {
		message.WriteString("\n\n")
		for _, detail := range details {
			message.WriteString(detail)
			message.WriteString("\n")
		}
	}

	return fmt.Errorf(message.String())
}

// GetProviderSuggestions returns provider-specific suggestions for errors
func GetProviderSuggestions(provider string) []string {
	switch provider {
	case "ollama":
		return []string{
			"Make sure Ollama is running with 'ollama serve'",
			"Check that the model is downloaded with 'ollama list'",
			"Try downloading the model using 'ollama pull <model>'",
			"Verify the endpoint URL (usually http://localhost:11434)",
		}
	case "openai":
		return []string{
			"Verify your API key is correct and not expired",
			"Check your OpenAI account for quota issues",
			"Make sure the model name is correct",
			"Verify the endpoint URL (usually https://api.openai.com)",
		}
	case "anthropic":
		return []string{
			"Verify your API key is correct and not expired",
			"Check your Anthropic account for quota issues",
			"Make sure the model name is correct",
			"Verify the endpoint URL (usually https://api.anthropic.com)",
		}
	default:
		return []string{
			"Check that the provider service is running",
			"Verify your credentials and configuration",
			"Make sure you have a stable internet connection",
		}
	}
}
