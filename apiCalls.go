package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Model represents the structure of each model in the API response
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Pricing     struct {
		Prompt     string `json:"prompt"`
		Completion string `json:"completion"`
	} `json:"pricing"`
}

// ModelsResponse represents the top-level API response
type ModelsResponse struct {
	Data []Model `json:"data"`
}

// ModelInfo holds the arrays of model information
type ModelInfo struct {
	IDs          []string
	Names        []string
	Descriptions []string
}

func FetchFreeModels() (*ModelInfo, error) {
	url := "https://openrouter.ai/api/v1/models"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var response ModelsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	info := &ModelInfo{
		IDs:          make([]string, 0),
		Names:        make([]string, 0),
		Descriptions: make([]string, 0),
	}

	for _, model := range response.Data {
		if model.Pricing.Prompt == "0" && model.Pricing.Completion == "0" {
			info.IDs = append(info.IDs, model.ID)
			info.Names = append(info.Names, model.Name)
			info.Descriptions = append(info.Descriptions, model.Description)
		}
	}
	println(info)
	return info, nil
}

func LlmCall(modelName string, usrInput string, history []Conversation, apiKey string) ([]byte, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	// Convert conversation history to Messages format
	messages := make([]Message, len(history)+1)
	for i, conv := range history {
		messages[i] = Message{
			Role: conv.Role,
			Content: []Content{
				{
					Type: "text",
					Text: conv.Content,
				},
			},
		}
	}

	// Add current user input
	messages = append(messages, Message{
		Role: "user",
		Content: []Content{
			{
				Type: "text",
				Text: usrInput,
			},
		},
	})

	// Create the request body
	reqBody := RequestBody{
		Model:    modelName,
		Messages: messages,
	}

	// Convert request body to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
