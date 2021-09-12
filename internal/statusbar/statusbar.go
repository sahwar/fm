package statusbar

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/knipferrc/fm/directory"
	"github.com/knipferrc/fm/formatter"
	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/internal/constants"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// Color is a struct that contains the foreground and background colors of the statusbar.
type Color struct {
	Background string
	Foreground string
}

// Model is a struct that contains all the properties of the statusbar.
type Model struct {
	Width              int
	Height             int
	TotalFiles         int
	Cursor             int
	TextInput          string
	ShowIcons          bool
	ShowCommandBar     bool
	InMoveMode         bool
	SelectedFile       fs.DirEntry
	ItemToMove         fs.DirEntry
	FirstColumnColors  Color
	SecondColumnColors Color
	ThirdColumnColors  Color
	FourthColumnColors Color
}

// NewModel creates an instance of a statusbar.
func NewModel(firstColumnColors, secondColumnColors, thirdColumnColors, fourthColumnColors Color) Model {
	return Model{
		Height:             1,
		TotalFiles:         0,
		Cursor:             0,
		TextInput:          "",
		ShowIcons:          true,
		ShowCommandBar:     false,
		InMoveMode:         false,
		SelectedFile:       nil,
		ItemToMove:         nil,
		FirstColumnColors:  firstColumnColors,
		SecondColumnColors: secondColumnColors,
		ThirdColumnColors:  thirdColumnColors,
		FourthColumnColors: fourthColumnColors,
	}
}

// ParseCommand parses the command and returns the command name and the arguments.
func ParseCommand(command string) (string, string) {
	// Split the command string into an array.
	cmdString := strings.Split(command, " ")

	// If theres only one item in the array, its a singular
	// command such as rm.
	if len(cmdString) == 1 {
		cmdName := cmdString[0]

		return cmdName, ""
	}

	// This command has two values, first one is the name
	// of the command, other is the value to pass back
	// to the UI to update.
	if len(cmdString) == 2 {
		cmdName := cmdString[0]
		cmdValue := cmdString[1]

		return cmdName, cmdValue
	}

	return "", ""
}

// GetHeight returns the height of the statusbar.
func (m Model) GetHeight() int {
	return m.Height
}

// SetContent sets the content of the statusbar.
func (m *Model) SetContent(totalFiles, cursor int, textInput string, showIcons, showCommandBar, inMoveMode bool, selectedFile, itemToMove fs.DirEntry) {
	m.TotalFiles = totalFiles
	m.Cursor = cursor
	m.TextInput = textInput
	m.ShowIcons = showIcons
	m.ShowCommandBar = showCommandBar
	m.InMoveMode = inMoveMode
	m.SelectedFile = selectedFile
	m.ItemToMove = itemToMove
}

// SetSize sets the size of the statusbar, useful when the terminal is resized.
func (m *Model) SetSize(width int) {
	m.Width = width
}

// Update updates the statusbar.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

}

// View returns a string representation of the statusbar.
func (m Model) View() string {
	width := lipgloss.Width
	logo := ""
	status := ""
	selectedFile := "N/A"
	fileCount := "0/0"

	if m.TotalFiles > 0 {
		if m.SelectedFile != nil {
			selectedFile = m.SelectedFile.Name()
			fileCount = fmt.Sprintf("%d/%d", m.Cursor+1, m.TotalFiles)

			currentPath, err := directory.GetWorkingDirectory()
			if err != nil {
				currentPath = constants.Directories.CurrentDirectory
			}

			fileInfo, err := m.SelectedFile.Info()
			if err != nil {
				return err.Error()
			}

			// Display some information about the currently seleted file including
			// its size, the mode and the current path.
			status = fmt.Sprintf("%s %s %s",
				formatter.ConvertBytesToSizeString(fileInfo.Size()),
				fileInfo.Mode().String(),
				currentPath,
			)
		}
	}

	if m.ShowCommandBar {
		status = m.TextInput
	}

	if m.ShowIcons {
		logo = fmt.Sprintf("%s %s", icons.IconDef["dir"].GetGlyph(), "FM")
	} else {
		logo = "FM"
	}

	selectedFileColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FirstColumnColors.Foreground)).
		Background(lipgloss.Color(m.FirstColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Render(truncate.StringWithTail(selectedFile, 30, "..."))

	fileCountColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.ThirdColumnColors.Foreground)).
		Background(lipgloss.Color(m.ThirdColumnColors.Background)).
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(m.Height).
		Render(fileCount)

	logoColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.FourthColumnColors.Foreground)).
		Background(lipgloss.Color(m.FourthColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Render(logo)

	statusColumn := lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.SecondColumnColors.Foreground)).
		Background(lipgloss.Color(m.SecondColumnColors.Background)).
		Padding(0, 1).
		Height(m.Height).
		Width(m.Width - width(selectedFileColumn) - width(fileCountColumn) - width(logoColumn)).
		Render(truncate.StringWithTail(status, uint(m.Width-width(selectedFileColumn)-width(fileCountColumn)-width(logoColumn)-3), "..."))

	return lipgloss.JoinHorizontal(lipgloss.Top,
		selectedFileColumn,
		statusColumn,
		fileCountColumn,
		logoColumn,
	)
}
