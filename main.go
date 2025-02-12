package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type     string    `json:"type,omitempty"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Conversation struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type model struct {
	choices  []string
	cursor   int
	selected string
}

type exitMenuModel struct {
	choices  []string
	cursor   int
	selected string
}

var models = []string{
	"google/gemini-2.0-flash-lite-preview-02-05:free",
	"google/gemini-2.0-flash-thinking-exp:free",
	"deepseek/deepseek-r1-distill-llama-70b:free",
	"meta-llama/llama-3.3-70b-instruct:free",
	"deepseek/deepseek-r1:free",
}

var exitMenuChoices = []string{
	"Continue chat",
	"Clear chat history",
	"Choose another model",
	"Exit program",
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("222"))
)

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func initialModel() model {
	return model{
		choices: models,
	}
}

func initialExitMenu() exitMenuModel {
	return exitMenuModel{
		choices: exitMenuChoices,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m exitMenuModel) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m exitMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	s := titleStyle.Render("Select an LLM:") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			choice = selectedStyle.Render(choice)
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n(use arrow keys to select, enter to confirm)\n"
	return s
}

func (m exitMenuModel) View() string {
	s := titleStyle.Render("What would you like to do?") + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
			choice = selectedStyle.Render(choice)
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n(use arrow keys to select, enter to confirm)\n"
	return s
}

func llmCall(modelName string, usrInput string, history []Conversation) ([]byte, error) {
	url := "https://openrouter.ai/api/v1/chat/completions"

	// Try to get API key from environment
	apiKey := os.Getenv("OPENROUTER_API_KEY")

	// If API key not found in environment, try loading from .env
	if apiKey == "" {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %v", err)
		}
		apiKey = os.Getenv("OPENROUTER_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENROUTER_API_KEY not found in environment or .env file")
		}
	}

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

func main() {
	history := make([]Conversation, 0)
	for {
		// Run the Bubble Tea program for model selection
		p := tea.NewProgram(initialModel())
		m, err := p.Run()
		if err != nil {
			fmt.Printf("Error running program: %v", err)
			os.Exit(1)
		}

		// Get the selected model
		modelName := m.(model).selected
		if modelName == "" {
			fmt.Println("No model selected")
			return
		}
		clearScreen()

		// Initialize conversation history

		// Chat loop
		scanner := bufio.NewScanner(os.Stdin)
		for {
			println("\nEnter your message (or /menu):")
			print("User: ")
			if !scanner.Scan() {
				break
			}

			input := scanner.Text()
			if input == "/menu" {
				clearScreen()
				// Show exit menu
				p := tea.NewProgram(initialExitMenu())
				m, err := p.Run()
				if err != nil {
					fmt.Printf("Error running program: %v", err)
					os.Exit(1)
				}

				choice := m.(exitMenuModel).selected
				if choice == "Exit program" {
					clearScreen()
					fmt.Println("Goodbye!")
					return
				} else if choice == "Choose another model" {
					clearScreen()
					break
				} else if choice == "Clear chat history" {
					history = make([]Conversation, 0)
					clearScreen()
					fmt.Println("Chat history cleared!")
					continue
				} else if choice == "Continue chat" {
					clearScreen()
					continue
				}
			}

			// Call LLM with history
			response, err := llmCall(modelName, input, history)
			if err != nil {
				panic(err)
			}

			// Parse the JSON response
			var result struct {
				Choices []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			}

			if err := json.Unmarshal(response, &result); err != nil {
				panic(err)
			}

			// Store the conversation
			history = append(history, Conversation{
				Role:    "user",
				Content: input,
			})

			// Get and store the response
			if len(result.Choices) > 0 {
				assistantResponse := result.Choices[0].Message.Content
				println("\nResponse:", assistantResponse)

				history = append(history, Conversation{
					Role:    "assistant",
					Content: assistantResponse,
				})
			}
		}
	}
}
