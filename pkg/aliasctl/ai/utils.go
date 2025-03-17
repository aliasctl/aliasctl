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

// ValidateEndpoint checks if the endpoint URL is valid
func ValidateEndpoint(endpoint string) error {
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		return fmt.Errorf("invalid endpoint URL. Must start with http:// or https://")
	}
	return nil
}

// GenerationPrompt creates a standardized prompt for alias generation
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

// ConversionPrompt creates a standardized prompt for alias conversion
func ConversionPrompt(alias, fromShell, toShell string) string {
	return fmt.Sprintf("Convert the following command from %s shell to %s shell. Provide only the final command without explanation: %s",
		fromShell, toShell, alias)
}

// ExtractAliasDefinition tries to extract the actual alias definition from response text
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

// MakeAPIRequest makes a generic API request with error handling
func MakeAPIRequest(method, url string, headers map[string]string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
