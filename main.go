package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
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

func extractApiKey() (string, error) {
	// Try to get API key from environment
	apiKey := os.Getenv("OPENROUTER_API_KEY")

	// If API key not found in environment, try loading from .env
	if apiKey == "" {
		if err := godotenv.Load(); err != nil {
			return "", fmt.Errorf("error loading .env file: %v", err)
		}
		apiKey = os.Getenv("OPENROUTER_API_KEY")
		if apiKey == "" {
			return "", fmt.Errorf("OPENROUTER_API_KEY not found in environment or .env file")
		}
	}
	return apiKey, nil
}

func main() {
	history := make([]Conversation, 0)
	apiKey, err := extractApiKey()
	if err != nil {
		// Handle error
		panic(err)
	}
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

		// Initialize spinner
		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		s.Color("magenta")
		s.Prefix = "Thinking... "

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

			// Start spinner
			fmt.Print("\nThinking... ")
			s.Start()

			// Call LLM with history
			response, err := LlmCall(modelName, input, history, apiKey)

			// Stop spinner
			s.Stop()
			fmt.Print("\r") // Clear the spinner line

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
				if !strings.HasSuffix(assistantResponse, "\n") {
					println("Response: " + assistantResponse)
				} else {
					print("Response: " + assistantResponse)
				}

				history = append(history, Conversation{
					Role:    "assistant",
					Content: assistantResponse,
				})
			}
		}
	}
}
