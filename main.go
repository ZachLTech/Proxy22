package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Location struct {
	Port     int    `json:"port"`
	Hostname string `json:"hostname"`
}

type Config struct {
	Locations map[string]Location `json:"locations"`
}

type Model struct {
	textInput   textinput.Model
	config      Config
	err         error
	destination string
}

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("error decoding config: %v", err)
	}
	return config, nil
}

func initialModel() Model {
	config, err := loadConfig("config.json")
	ti := textinput.New()
	ti.Placeholder = "Enter full domain (e.g., exit.zachl.tech)"
	ti.Focus()
	return Model{
		textInput: ti,
		config:    config,
		err:       err,
	}
}
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.destination = m.textInput.Value()
			if location, exists := m.config.Locations[m.destination]; exists {
				m.textInput.Placeholder = fmt.Sprintf("Routing to %s on port %d (hostname: %s)\n", m.destination, location.Port, location.Hostname)
				m.textInput.SetValue("")
				return m, tea.Quit
			} else {
				m.textInput.Placeholder = fmt.Sprintf("no destination found at location: %s", m.destination)
				m.textInput.SetValue("")
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\n%s\n\n(esc to quit)", m.err, m.textInput.View())
	}
	return fmt.Sprintf("MAKE THIS TITLE CUSTOMIZABLE WITH CONFIG!\n\nWhere are you connecting? Enter the full domain (including subdomains if applicable):\n\n%s\n\n(esc to quit)", m.textInput.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
