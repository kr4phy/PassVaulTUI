package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	stateEnterPassword     = 0
	statePasswordsList     = 1
	statePasswordDetail    = 2
	stateAddEntry          = 3
	stateChangeMaster      = 4
	statePasswordGenerator = 5
)

var (
	windowStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1) // 좌우 여백
	formWindowStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1)
	titleStyle   = lipgloss.NewStyle().Bold(true).AlignHorizontal(lipgloss.Center)
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	textStyle    = lipgloss.NewStyle()
)

type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	vpWidth          int
	vpHeight         int
	currentState     int
	storeExists      bool
	passwordInput    textinput.Model
	masterPass       string
	passStorage      [][3]string
	passList         list.Model
	isEditing        bool
	editIndex        int
	addTitleInput    textinput.Model
	addIDInput       textinput.Model
	addPassInput     textinput.Model
	addFocus         int
	changePassInput  textinput.Model
	genLengthInput   textinput.Model
	genUseUpper      bool
	genUseNumber     bool
	genUseSpecial    bool
	genFocus         int
	generatedPass    string
	chosenCredential list.Item
	Ticks            int
	err              error
}

func initialModel() model {
	pInput := textinput.New()
	pInput.Placeholder = "Enter master password"
	pInput.Focus()
	pInput.CharLimit = 156
	pInput.Width = 24
	pInput.EchoMode = textinput.EchoPassword
	pInput.EchoCharacter = '*'
	passList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	passList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new entry")),
			key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit entry")),
			key.NewBinding(key.WithKeys("d", "backspace", "delete"), key.WithHelp("d/backspace/delete", "remove")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open detail")),
			key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "change master")),
			key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "generate password")),
		}
	}
	var storeExists bool
	if _, err := os.Stat("data.bin"); err != nil {
		storeExists = false
	} else {
		storeExists = true
	}
	addTitle := textinput.New()
	addTitle.Placeholder = "title"
	addTitle.CharLimit = 156
	addTitle.Width = 24
	addID := textinput.New()
	addID.Placeholder = "id"
	addID.CharLimit = 156
	addID.Width = 24
	addPass := textinput.New()
	addPass.Placeholder = "password"
	addPass.CharLimit = 156
	addPass.Width = 24
	addPass.EchoMode = textinput.EchoPassword
	addPass.EchoCharacter = '*'
	changePass := textinput.New()
	changePass.Placeholder = "Enter new master password"
	changePass.CharLimit = 156
	changePass.Width = 24
	changePass.EchoMode = textinput.EchoPassword
	changePass.EchoCharacter = '*'
	genLength := textinput.New()
	genLength.Placeholder = "length"
	genLength.CharLimit = 4
	genLength.Width = 6
	genLength.SetValue("16")
	return model{
		vpWidth:          0,
		vpHeight:         0,
		currentState:     0,
		storeExists:      storeExists,
		passwordInput:    pInput,
		masterPass:       "",
		passStorage:      nil,
		passList:         passList,
		isEditing:        false,
		editIndex:        -1,
		addTitleInput:    addTitle,
		addIDInput:       addID,
		addPassInput:     addPass,
		addFocus:         0,
		changePassInput:  changePass,
		genLengthInput:   genLength,
		genUseUpper:      true,
		genUseNumber:     true,
		genUseSpecial:    true,
		genFocus:         0,
		generatedPass:    "",
		chosenCredential: item{},
		Ticks:            0,
		err:              nil,
	}
}

func (m model) Init() tea.Cmd { return tea.Batch(textinput.Blink, tick()) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//var cmd tea.Cmd
	if m.Ticks == 5 {
		m.Ticks = 0
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := passListStyle.GetFrameSize()
		m.passList.SetSize(msg.Width-h*2, msg.Height-v*2)
		m.vpWidth = msg.Width
		m.vpHeight = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tickMsg:
		m.Ticks++
		return m, tick()
	}
	switch m.currentState {
	case stateEnterPassword:
		return UpdateEnterPassword(msg, m)
	case statePasswordsList:
		return UpdatePasswordList(msg, m)
	case statePasswordDetail:
		return UpdatePasswordDetail(msg, m)
	case stateAddEntry:
		return UpdateEditEntry(msg, m)
	case stateChangeMaster:
		return UpdateChangeMasterPass(msg, m)
	case statePasswordGenerator:
		return UpdatePasswordGenerator(msg, m)
	default:
		return UpdateEnterPassword(msg, m)
	}
}

func (m model) View() string {
	var content string
	switch m.currentState {
	case stateEnterPassword:
		content = EnterPasswordView(m)
	case statePasswordsList:
		content = PasswordListView(m)
	case statePasswordDetail:
		content = PasswordDetailView(m)
	case stateAddEntry:
		content = EditEntryView(m)
	case stateChangeMaster:
		content = ChangeMasterPassView(m)
	case statePasswordGenerator:
		content = PasswordGeneratorView(m)
	default:
		content = EnterPasswordView(m)
	}
	return windowStyle.Render(content)
}
