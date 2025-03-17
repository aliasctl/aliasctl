package ai

import (
	"encoding/json"
	"fmt"
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

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  op.Model,
		"prompt": prompt,
		"stream": false, // Disable streaming to get full response
	})
	if err != nil {
		return "", err
	}

	respBody, err := MakeAPIRequest("POST", op.Endpoint+"/api/generate", nil, requestBody)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %v", err)
	}

	// Define a struct to match the Ollama response
	var ollamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	err = json.Unmarshal(respBody, &ollamaResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing Ollama response: %v. Raw response: %s", err, string(respBody))
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

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  op.Model,
		"prompt": prompt,
		"stream": false, // Disable streaming to get full response
	})
	if err != nil {
		return "", err
	}

	respBody, err := MakeAPIRequest("POST", op.Endpoint+"/api/generate", nil, requestBody)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %v", err)
	}

	// Define a struct to match the Ollama response
	var ollamaResponse struct {
		Response string `json:"response"`
		Done     bool   `json:"done"`
	}

	err = json.Unmarshal(respBody, &ollamaResponse)
	if err != nil {
		return "", fmt.Errorf("error parsing Ollama response: %v. Raw response: %s", err, string(respBody))
	}

	return ExtractAliasDefinition(ollamaResponse.Response), nil
}
