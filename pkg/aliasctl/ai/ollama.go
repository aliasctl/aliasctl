package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OllamaProvider implements Provider for Ollama.
type OllamaProvider struct {
	Endpoint string // The Ollama endpoint URL
	Model    string // The Ollama model name
}

// GenerateAlias generates an alias using Ollama AI
func (op *OllamaProvider) GenerateAlias(command, shellType string) (string, error) {
	if err := ValidateEndpoint(op.Endpoint); err != nil {
		return "", err
	}

	prompt := GenerationPrompt(command, shellType)

	requestBody, err := json.Marshal(map[string]any{
		"model":  op.Model,
		"prompt": prompt,
		"stream": false, // Disable streaming to get full response
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Ollama request: %w", err)
	}

	respBody, err := MakeAPIRequest("POST", op.Endpoint+"/api/generate", nil, requestBody)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return "", fmt.Errorf("failed to connect to Ollama at %s: make sure Ollama is running with 'ollama serve'", op.Endpoint)
		}
		return "", fmt.Errorf("ollama request failed: %w", err)
	}

	// Define a struct to match the Ollama response
	var ollamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
		Error    string `json:"error"`
	}

	err = json.Unmarshal(respBody, &ollamaResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing ollama response: %w\n\nRaw response: %s", err, limitResponseText(string(respBody), 200))
	}

	// Check if Ollama returned an error
	if ollamaResponse.Error != "" {
		// Check for model-related errors
		if strings.Contains(ollamaResponse.Error, "model") && strings.Contains(ollamaResponse.Error, "not found") {
			return "", fmt.Errorf("ollama model '%s' not found: run 'ollama pull %s' to download it first", op.Model, op.Model)
		}
		return "", fmt.Errorf("ollama error: %s", ollamaResponse.Error)
	}

	// Parse the alias from the response
	return ExtractAliasDefinition(ollamaResponse.Response), nil
}

// ConvertAlias converts an alias using the Ollama AI service.
func (op *OllamaProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	if err := ValidateEndpoint(op.Endpoint); err != nil {
		return "", err
	}

	prompt := ConversionPrompt(alias, fromShell, toShell)

	requestBody, err := json.Marshal(map[string]any{
		"model":  op.Model,
		"prompt": prompt,
		"stream": false, // Disable streaming to get full response
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Ollama request: %w", err)
	}

	respBody, err := MakeAPIRequest("POST", op.Endpoint+"/api/generate", nil, requestBody)
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			return "", fmt.Errorf("failed to connect to Ollama at %s: make sure Ollama is running with 'ollama serve'", op.Endpoint)
		}
		return "", fmt.Errorf("ollama request failed: %w", err)
	}

	// Define a struct to match the Ollama response
	var ollamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
		Error    string `json:"error"`
	}

	err = json.Unmarshal(respBody, &ollamaResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing ollama response: %w\n\nRaw response: %s", err, limitResponseText(string(respBody), 200))
	}

	// Check if Ollama returned an error
	if ollamaResponse.Error != "" {
		// Check for model-related errors
		if strings.Contains(ollamaResponse.Error, "model") && strings.Contains(ollamaResponse.Error, "not found") {
			return "", fmt.Errorf("ollama model '%s' not found: run 'ollama pull %s' to download it first", op.Model, op.Model)
		}
		return "", fmt.Errorf("ollama error: %s", ollamaResponse.Error)
	}

	return ExtractAliasDefinition(ollamaResponse.Response), nil
}
