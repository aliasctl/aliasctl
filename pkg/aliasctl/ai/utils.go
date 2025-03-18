package ai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Common timeout for all API requests
const defaultTimeout = 30 * time.Second

// Common HTTP client that can be reused
var httpClient = &http.Client{
	Timeout: defaultTimeout,
}

// ValidateEndpoint checks if the endpoint URL is valid.
// It ensures the URL starts with http:// or https://.
// Returns an error if the URL format is invalid.
func ValidateEndpoint(endpoint string) error {
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return fmt.Errorf("invalid endpoint URL '%s': must start with http:// or https://", endpoint)
	}
	return nil
}

// GenerationPrompt creates a standardized prompt for alias generation.
// It formats a prompt for AI models to create shell aliases based on the command and shell type.
// The prompt includes context, requirements, and formatting instructions.
func GenerationPrompt(command, shellType string) string {
	return fmt.Sprintf(`You are a shell alias creation expert for %s shell.

Task: Create a concise, memorable alias for the following command:
%s

Requirements:
- The alias name should be short but descriptive
- Follow standard naming conventions for %s aliases
- The alias should be intuitive and easy to remember
- Don't abbreviate too aggressively, though initials like kgp for kubectl get pods are acceptable.
- Avoid using special characters or spaces in the alias
- Ensure the alias is unique and doesn't conflict with existing commands in the shell
- Consider common aliases in the %s ecosystem

Response format:
Provide ONLY the complete alias definition in the correct syntax for %s shell.
- For bash/zsh: alias name='command'
- For PowerShell: Set-Alias name command or function name { command }
- For CMD: doskey name=command
- For fish: alias name 'command' or function name\n    command\nend

Do not include any explanations, preambles, or additional text.`,
		shellType, command, shellType, shellType, shellType)
}

// ConversionPrompt creates a standardized prompt for alias conversion.
// It formats a prompt for AI models to convert an alias between different shell formats.
// The prompt specifies the source and target shell types and the alias to convert.
func ConversionPrompt(alias, fromShell, toShell string) string {
	return fmt.Sprintf("Convert the following command from %s shell to %s shell. Provide only the final command without explanation: %s",
		fromShell, toShell, alias)
}

// ExtractAliasDefinition tries to extract the actual alias definition from response text.
// It parses AI-generated responses to find the valid alias definition, looking for
// common patterns like "alias", "function", "Set-Alias", or "doskey" prefixes.
// Returns the entire content if no specific pattern is found.
func ExtractAliasDefinition(content string) string {
	// Trim any leading/trailing whitespace
	content = strings.TrimSpace(content)

	// Extract just the command if possible
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "alias ") ||
			strings.HasPrefix(line, "function ") ||
			strings.HasPrefix(line, "Set-Alias ") ||
			strings.HasPrefix(line, "doskey ") {
			return line
		}
	}

	return content
}

// MakeAPIRequest makes a generic API request with error handling.
// It creates an HTTP request with the specified method, URL, headers, and body,
// then executes it and processes the response.
// Returns the response body and any error encountered during the request.
// Provides detailed error messages based on HTTP status codes and common error patterns.
func MakeAPIRequest(method, url string, headers map[string]string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request to %s: %w", url, err)
	}

	// Set default content type
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Execute the request
	resp, err := httpClient.Do(req)
	if err != nil {
		// Check for common network errors and provide better messages
		if strings.Contains(err.Error(), "connection refused") {
			return nil, fmt.Errorf("connection refused to %s: service may not be running or the endpoint is incorrect", url)
		}
		if strings.Contains(err.Error(), "no such host") {
			return nil, fmt.Errorf("host not found for %s: check your network connection and the endpoint URL", url)
		}
		if strings.Contains(err.Error(), "timeout") {
			return nil, fmt.Errorf("request to %s timed out after %s: the service might be overloaded or unreachable", url, defaultTimeout)
		}
		return nil, fmt.Errorf("error connecting to %s: %w", url, err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response from %s: %w", url, err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		// Attempt to provide more context based on status code
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return respBody, fmt.Errorf("API authentication error (status 401): invalid or missing API key")
		case http.StatusForbidden:
			return respBody, fmt.Errorf("API authorization error (status 403): your API key doesn't have permission for this operation")
		case http.StatusNotFound:
			return respBody, fmt.Errorf("API resource not found (status 404): the endpoint URL or API version might be incorrect")
		case http.StatusTooManyRequests:
			return respBody, fmt.Errorf("API rate limit exceeded (status 429): try again later or check your API usage limits")
		case http.StatusInternalServerError:
			return respBody, fmt.Errorf("API server error (status 500): the service might be experiencing issues")
		default:
			return respBody, fmt.Errorf("API error (status %d): %s", resp.StatusCode, limitResponseText(string(respBody), 200))
		}
	}

	return respBody, nil
}

// limitResponseText limits response text to a maximum length.
// It truncates text that exceeds the maximum length and adds an ellipsis.
// Used to prevent overly verbose error messages when API responses are large.
func limitResponseText(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength] + "..."
}
