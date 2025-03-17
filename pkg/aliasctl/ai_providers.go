package aliasctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ConvertAlias converts an alias from one shell to another using the configured AI provider.
func (am *AliasManager) ConvertAlias(name, targetShell string, providerName string) (string, error) {
	if !am.AIConfigured {
		return "", fmt.Errorf("AI provider not configured")
	}

	commands, exists := am.Aliases[name]
	if !exists {
		return "", fmt.Errorf("alias '%s' not found", name)
	}

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

	if command == "" {
		return "", fmt.Errorf("command for shell '%s' not found", am.Shell)
	}

	var provider AIProvider
	if providerName == "" {
		// Default to the legacy AIProvider if no provider specified
		provider = am.AIProvider
	} else {
		var exists bool
		provider, exists = am.AIProviders[providerName]
		if !exists {
			available := am.GetAvailableProviders()
			return "", fmt.Errorf("AI provider '%s' not configured. Available providers: %v",
				providerName, strings.Join(available, ", "))
		}
	}

	return provider.ConvertAlias(command, string(am.Shell), targetShell)
}

// GenerateAlias generates an alias suggestion for the given command using AI.
func (am *AliasManager) GenerateAlias(command string, providerName string) (string, error) {
	if !am.AIConfigured {
		return "", fmt.Errorf("AI provider not configured")
	}

	var provider AIProvider
	if providerName == "" {
		// Default to the legacy AIProvider if no provider specified
		provider = am.AIProvider
	} else {
		var exists bool
		provider, exists = am.AIProviders[providerName]
		if !exists {
			available := am.GetAvailableProviders()
			return "", fmt.Errorf("AI provider '%s' not configured. Available providers: %v",
				providerName, strings.Join(available, ", "))
		}
	}

	switch p := provider.(type) {
	case *OllamaProvider:
		return generateAliasOllama(p, command, string(am.Shell))
	case *OpenAIProvider:
		return generateAliasOpenAI(p, command, string(am.Shell))
	case *AnthropicProvider:
		return generateAliasAnthropic(p, command, string(am.Shell))
	default:
		return "", fmt.Errorf("unsupported AI provider type")
	}
}

// generateAliasOllama generates an alias using the Ollama AI service.
func generateAliasOllama(op *OllamaProvider, command, shellType string) (string, error) {
	// Validate endpoint
	if !strings.HasPrefix(op.Endpoint, "http://") && !strings.HasPrefix(op.Endpoint, "https://") {
		return "", fmt.Errorf("invalid Ollama endpoint. Must start with http:// or https://")
	}

	prompt := fmt.Sprintf("I have the following shell command for %s shell: %s\nSuggest a concise and memorable alias for this command. Respond only with the full alias definition including the name and command.", shellType, command)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  op.Model,
		"prompt": prompt,
		"stream": false, // Disable streaming to get full response
	})
	if err != nil {
		return "", err
	}

	// Add timeout to prevent hanging
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Post(op.Endpoint+"/api/generate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error connecting to Ollama: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Define a struct to match the Ollama response
	var ollamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing Ollama response: %v. Raw response: %s", err, string(body))
	}

	// Parse the alias from the response
	return ParseAliasCommand(ollamaResponse.Response), nil
}

// generateAliasOpenAI generates an alias using the OpenAI-compatible API.
func generateAliasOpenAI(op *OpenAIProvider, command, shellType string) (string, error) {
	// Validate endpoint
	if !strings.HasPrefix(op.Endpoint, "http://") && !strings.HasPrefix(op.Endpoint, "https://") {
		return "", fmt.Errorf("invalid OpenAI endpoint. Must start with http:// or https://")
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": op.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a shell alias creation assistant. Given a command, suggest a concise and memorable alias. Respond only with the full alias definition including the name and command, appropriate for the specified shell.",
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Create an alias for this %s shell command: %s", shellType, command),
			},
		},
		"temperature": 0.3, // Moderate creativity
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", op.Endpoint+"/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+op.APIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return ParseAliasCommand(content), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from OpenAI")
}

// generateAliasAnthropic generates an alias using the Anthropic Claude API.
func generateAliasAnthropic(ap *AnthropicProvider, command, shellType string) (string, error) {
	// Validate endpoint
	if !strings.HasPrefix(ap.Endpoint, "http://") && !strings.HasPrefix(ap.Endpoint, "https://") {
		return "", fmt.Errorf("invalid Anthropic endpoint. Must start with http:// or https://")
	}

	// Build the request payload
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": ap.Model,
		"messages": []map[string]string{
			{
				"role": "user",
				"content": fmt.Sprintf("Create a concise and memorable shell alias for this %s command: %s\n\nRespond only with the full alias definition including the name and command.",
					shellType, command),
			},
		},
		"max_tokens":  300,
		"temperature": 0.3, // Moderate creativity
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ap.Endpoint+"/v1/messages", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", ap.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01") // Use appropriate API version

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anthropic API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Process the response content
	var responseText string
	for _, content := range result.Content {
		if content.Type == "text" {
			responseText = content.Text
			break
		}
	}

	if responseText == "" {
		return "", fmt.Errorf("no text response found in Anthropic Claude reply")
	}

	return ParseAliasCommand(responseText), nil
}

// ConvertAlias for OllamaProvider converts an alias using the Ollama AI service.
func (op *OllamaProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	// Validate endpoint
	if !strings.HasPrefix(op.Endpoint, "http://") && !strings.HasPrefix(op.Endpoint, "https://") {
		return "", fmt.Errorf("invalid Ollama endpoint. Must start with http:// or https://")
	}

	prompt := fmt.Sprintf("Convert the following command from %s shell to %s shell, providing only the final command without explanation: %s", fromShell, toShell, alias)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  op.Model,
		"prompt": prompt,
		"stream": false, // Disable streaming to get full response
	})
	if err != nil {
		return "", err
	}

	// Add timeout to prevent hanging
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Post(op.Endpoint+"/api/generate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error connecting to Ollama: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// Define a struct to match the Ollama response
	var ollamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	err = json.Unmarshal(body, &ollamaResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing Ollama response: %v. Raw response: %s", err, string(body))
	}

	// Trim any leading/trailing whitespace
	convertedAlias := strings.TrimSpace(ollamaResponse.Response)

	// Extract the actual alias definition if possible
	aliasLines := strings.Split(convertedAlias, "\n")
	for _, line := range aliasLines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "alias ") ||
			strings.HasPrefix(line, "function ") ||
			strings.HasPrefix(line, "Set-Alias ") ||
			strings.HasPrefix(line, "doskey ") {
			return line, nil
		}
	}

	return convertedAlias, nil
}

// ConvertAlias for OpenAIProvider converts an alias using the OpenAI-compatible API.
func (op *OpenAIProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	// Validate endpoint
	if !strings.HasPrefix(op.Endpoint, "http://") && !strings.HasPrefix(op.Endpoint, "https://") {
		return "", fmt.Errorf("invalid OpenAI endpoint. Must start with http:// or https://")
	}

	prompt := fmt.Sprintf("Convert the following command from %s shell to %s shell. Provide only the final command without explanation: %s", fromShell, toShell, alias)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": op.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a utility that converts command line aliases between different shells. Respond only with the converted command, no explanation.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.2, // Lower temperature for more deterministic results
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", op.Endpoint+"/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+op.APIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					// Clean up the response
					content = strings.TrimSpace(content)

					// Extract just the command if possible
					lines := strings.Split(content, "\n")
					for _, line := range lines {
						line = strings.TrimSpace(line)
						if strings.HasPrefix(line, "alias ") ||
							strings.HasPrefix(line, "function ") ||
							strings.HasPrefix(line, "Set-Alias ") ||
							strings.HasPrefix(line, "doskey ") {
							return line, nil
						}
					}

					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from OpenAI")
}

// ConvertAlias for AnthropicProvider converts an alias using the Anthropic Claude API.
func (ap *AnthropicProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	// Validate endpoint
	if !strings.HasPrefix(ap.Endpoint, "http://") && !strings.HasPrefix(ap.Endpoint, "https://") {
		return "", fmt.Errorf("invalid Anthropic endpoint. Must start with http:// or https://")
	}

	prompt := fmt.Sprintf("Convert the following command from %s shell to %s shell. Provide only the final command without explanation: %s", fromShell, toShell, alias)

	// Build the request payload
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": ap.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  300,
		"temperature": 0.1, // Lower temperature for more deterministic results
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ap.Endpoint+"/v1/messages", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", ap.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01") // Use appropriate API version

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anthropic API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Process the response content
	var responseText string
	for _, content := range result.Content {
		if content.Type == "text" {
			responseText = content.Text
			break
		}
	}

	if responseText == "" {
		return "", fmt.Errorf("no text response found in Anthropic Claude reply")
	}

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
			return line, nil
		}
	}

	return responseText, nil
}
