package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"

	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mgale/bash-history.git/internal/defaults"
	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/xorcare/pointer"
	"golang.org/x/term"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(0)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	//quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	tokenSeparators   = map[rune]bool{'-': true, '/': true, '_': true, '=': true}
	defaultListHeight = 40
	defaultListWidth  = 200
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type searchTUI struct {
	resultsList list.Model
	searchBox   textinput.Model
	//Used to hold the content from the search box while we are searching
	inputChars []rune
	username   string
	//The selected item from the results list
	choice        string
	searchCounter int
}

func (s searchTUI) Init() tea.Cmd {
	return nil
}

func (s searchTUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	runSearch := false
	runSearchForce := false
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		updateWindowsSize(msg.Width, msg.Height)
		return s, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit
		case tea.KeyEsc:
			return s, tea.Quit
		case tea.KeyEnter:
			i, ok := s.resultsList.SelectedItem().(item)
			if ok {
				s.choice = string(i)
				cmdStr = s.choice
			}
			return s, tea.Quit
		case tea.KeyRunes:
			log.Println("KeyRunes: ", msg.Runes)
			s.inputChars = append(s.inputChars, msg.Runes...)
			runSearch = true
		case tea.KeySpace:
			s.inputChars = append(s.inputChars, msg.Runes...)
			runSearch = true
		case tea.KeyBackspace:
			if len(s.inputChars) > 0 {
				s.inputChars = s.inputChars[:len(s.inputChars)-1]
				runSearch = true
			}
			// Provides a refresh of the results list
			runSearchForce = true
		case tea.KeyCtrlR:
			runSearchForce = true
		case tea.KeyCtrlS:
			runSearchForce = true
		}
	}

	s.searchBox.SetValue(string(s.inputChars))

	// Only execute a search and update when we have at least 3 characters and
	// have received a textbox update.
	if (runSearch && len(s.inputChars) > 2) || runSearchForce {
		var err error
		s.searchCounter++
		s.resultsList, err = createItemsListModel(tsclient, s.username, s.inputChars, s.searchCounter)
		if err != nil {
			log.Fatal(err)
		}
	}

	var cmd tea.Cmd
	s.resultsList, _ = s.resultsList.Update(msg)
	s.searchBox, cmd = s.searchBox.Update(msg)

	return s, cmd
}

func (s searchTUI) View() string {
	return "\n" + s.resultsList.View() + "\n" + s.searchBox.View()
}

func initialModel(defaultQuery []rune) searchTUI {

	_ = setDefaultWindowSize()
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	items, err := createItemsListModel(tsclient, currentUser.Username, defaultQuery, 1)
	if err != nil {
		log.Fatalf(err.Error())
	}

	t := textinput.NewModel()
	t.CharLimit = 1024
	t.Placeholder = "Search..."
	t.Prompt = "Search: "

	return searchTUI{
		resultsList:   items,
		searchBox:     t,
		username:      currentUser.Username,
		searchCounter: 1,
	}
}

func createItemsListModel(tsclient *typesense.Client, username string, query []rune, searchCounter int) (list.Model, error) {
	log.Println("Creating items list model")
	results, err := getHistoryEvents(tsclient, username, query)
	if err != nil {
		return list.Model{}, err
	}

	items := createItemsFromSearchResult(results)

	l := list.New(items, itemDelegate{}, defaultListWidth, defaultListHeight)
	l.Title = fmt.Sprintf("Command history, showing %d out of possible %d, RT %dms, TSC: %d", len(items), *results.Found, *results.SearchTimeMs, searchCounter)
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.KeyMap = createCustomKeyBinds()

	return l, nil
}

// createItemsFromSearchResult creates a list of items from the search result
// provides a way to convert the search result into a list of items after sanitizing the data
func createItemsFromSearchResult(results *api.SearchResult) []list.Item {
	items := make([]list.Item, 0, len(*results.Hits))
	for _, hit := range *results.Hits {
		kv := hit.Document
		cmdString := (*kv)["command"].(string)
		cmdString = strings.TrimSpace(cmdString)
		if len(cmdString) > 0 {
			items = append(items, item(cmdString))
		}
	}
	return items
}

// getHistoryEvents returns a list of history events from typesense
func getHistoryEvents(tsclient *typesense.Client, user string, query []rune) (*api.SearchResult, error) {
	log.Println("Getting history events")
	queryString := createQueryString(query)
	searchParameters := &api.SearchCollectionParams{
		Q:        queryString,
		QueryBy:  "command",
		Infix:    pointer.String("off"),
		SortBy:   pointer.String("timestamp:desc"),
		PerPage:  pointer.Int(250),
		NumTypos: pointer.Int(3),
	}

	if user != "" {
		searchParameters.FilterBy = pointer.String(fmt.Sprintf("username:%s", user))
	}

	//log.Printf("Querying typesense with: %s\n", searchParameters)
	results, err := tsclient.Collection(defaults.CollectionName).Documents().Search(searchParameters)
	if err != nil {
		return nil, err
	}

	return results, nil
}

/*
A custom key binding is used because the default list navigation keys contain letters like
h,l,j, etc to navigate left, right, down, etc. This is not ideal because we want those keys
to be part of the search query.
From: https://github.com/charmbracelet/bubbles/blob/db06ae17d341503687ca5bdebcf0e8d6d90a9c86/list/keys.go#L34
*/
func createCustomKeyBinds() list.KeyMap {
	return list.KeyMap{
		// Browsing.
		CursorUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("left", "pgup"),
			key.WithHelp("←/pgup", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right", "pgdown"),
			key.WithHelp("→/pgdn", "next page"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "go to start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "go to end"),
		),

		// Toggle help.
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		// Quitting.
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
}

func createQueryString(input []rune) string {
	result := []rune{}
	for _, c := range input {
		// if the character is a token separator, skip it.
		if _, ok := tokenSeparators[c]; ok {
			continue
		}
		result = append(result, c)
	}

	log.Println("Query string created:", string(result))
	return string(result)
}

// Set / Update the window size
func updateWindowsSize(width, height int) {
	// This needs to account for the border and padding of the list
	defaultListWidth = width - 5
	// This needs to account for the title, border, padding, and status bar of the list
	// Plus the text input box at the bottom
	defaultListHeight = height - 8
	log.Println("Set window size to", defaultListWidth, defaultListHeight)
}

func setDefaultWindowSize() error {
	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		return err
	}

	updateWindowsSize(width, height)
	return nil
}
