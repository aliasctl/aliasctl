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
func (am *AliasManager) ConvertAlias(name, targetShell string) (string, error) {
	if !am.AIConfigured {
		return "", fmt.Errorf("AI provider not configured")
	}

	command, exists := am.Aliases[name]
	if !exists {
		return "", fmt.Errorf("alias '%s' not found", name)
	}

	return am.AIProvider.ConvertAlias(command, string(am.Shell), targetShell)
}

// ConvertAlias for OllamaProvider converts an alias using the Ollama AI service.
func (op *OllamaProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	// Validate endpoint
	if !strings.HasPrefix(op.Endpoint, "http://") && !strings.HasPrefix(op.Endpoint, "https://") {
		return "", fmt.Errorf("invalid Ollama endpoint. Must start with http:// or https://")
	}

	prompt := fmt.Sprintf("Convert the following command from %s shell to %s shell: %s", fromShell, toShell, alias)

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

	// Extract the actual alias definition
	aliasLines := strings.Split(convertedAlias, "\n")
	for _, line := range aliasLines {
		if strings.HasPrefix(line, "alias ") {
			return strings.TrimSpace(line), nil
		}
	}

	return convertedAlias, nil
}

// ConvertAlias for OpenAIProvider converts an alias using the OpenAI-compatible API.
func (op *OpenAIProvider) ConvertAlias(alias, fromShell, toShell string) (string, error) {
	prompt := fmt.Sprintf("Convert the following command from %s shell to %s shell: %s", fromShell, toShell, alias)

	requestBody, err := json.Marshal(map[string]interface{}{
		"model": op.Model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an assistant that converts command line aliases between different shells.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", op.Endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+op.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unexpected response format from OpenAI")
}
