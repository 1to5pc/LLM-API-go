package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

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
