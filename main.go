package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			row := m.table.SelectedRow()
			var cmdStr = []string{"build"}
			if len(name) == 0 {
				cmdStr = append(cmdStr, path)
			} else {
				// name = filepath.Join(exPath, name)
				cmdStr = append(cmdStr, "-o", name, path)
			}

			log.Println("go", strings.Join(cmdStr, " "))

			cmd := exec.Command("go", cmdStr...)
			cmd.Dir = exPath
			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", row[1]), fmt.Sprintf("GOARCH=%s", row[2]))

			b, err := cmd.CombinedOutput()
			if err != nil {
				log.Println(err)
				log.Fatal(string(b))
			}

			log.Println("executable generated. q to exit")
			return m, nil
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

var (
	exPath string
	path   string
	name   string
)

func main() {
	flag.StringVar(&path, "path", "main.go", "path to main file")
	flag.StringVar(&name, "name", "", "name of executable. mimics go build behaviour")

	flag.Parse()

	// log.Println(exPath)

	dir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal("invalid path")
	}

	exPath = dir
	// path = filepath.Join(exPath, path)

	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "OS", Width: 5},
		{Title: "Arch", Width: 5},
	}

	// "amd64", "arm", "arm64"
	rows := []table.Row{
		{"1", "linux", "amd64"},
		{"2", "darwin", "amd64"},
		{"3", "linux", "arm"},
		{"4", "linux", "arm64"},
		{"5", "darwin", "arm"},
		{"6", "darwin", "arm64"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
