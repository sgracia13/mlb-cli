package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"mlb-cli/internal/api"
	"mlb-cli/internal/models"
)

// View represents the current view/screen
type View int

const (
	ViewTeams View = iota
	ViewRoster
	ViewPlayer
	ViewStandings
	ViewSchedule
)

// Tab represents the main navigation tabs
type Tab int

const (
	TabTeams Tab = iota
	TabStandings
	TabSchedule
)

// Model is the main TUI model
type Model struct {
	// API client
	client *api.Client

	// Current view state
	currentView View
	currentTab  Tab

	// Navigation history for back functionality
	history []View

	// Data
	teams       []models.Team
	roster      []models.RosterEntry
	player      *models.Player
	playerStats *models.PlayerStatsResponse
	standings   *models.StandingsResponse
	schedule    *models.ScheduleResponse

	// Selected items
	selectedTeam   *models.Team
	selectedPlayer *models.RosterEntry

	// UI state
	cursor       int
	loading      bool
	err          error
	filterMode   bool
	filterText   string
	filteredList []int // indices of filtered items

	// Components
	spinner   spinner.Model
	textInput textinput.Model
	keys      KeyMap

	// Window size
	width  int
	height int
}

// NewModel creates a new TUI model
func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	ti := textinput.New()
	ti.Placeholder = "Type to filter..."
	ti.CharLimit = 50
	ti.Width = 30

	return Model{
		client:      api.NewClient(),
		currentView: ViewTeams,
		currentTab:  TabTeams,
		history:     make([]View, 0),
		cursor:      0,
		loading:     true,
		spinner:     s,
		textInput:   ti,
		keys:        DefaultKeyMap,
		width:       80,
		height:      24,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadTeams,
	)
}

// Messages for async operations
type teamsLoadedMsg struct {
	teams []models.Team
	err   error
}

type rosterLoadedMsg struct {
	roster []models.RosterEntry
	err    error
}

type standingsLoadedMsg struct {
	standings *models.StandingsResponse
	err       error
}

type scheduleLoadedMsg struct {
	schedule *models.ScheduleResponse
	err      error
}

type playerStatsLoadedMsg struct {
	stats *models.PlayerStatsResponse
	err   error
}

// Command functions
func (m Model) loadTeams() tea.Msg {
	resp, err := m.client.GetTeams()
	if err != nil {
		return teamsLoadedMsg{err: err}
	}
	return teamsLoadedMsg{teams: resp.Teams}
}

func (m Model) loadRoster(teamID string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.GetRoster(teamID)
		if err != nil {
			return rosterLoadedMsg{err: err}
		}
		return rosterLoadedMsg{roster: resp.Roster}
	}
}

func (m Model) loadStandings() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.GetStandings("2024")
		if err != nil {
			return standingsLoadedMsg{err: err}
		}
		return standingsLoadedMsg{standings: resp}
	}
}

func (m Model) loadSchedule() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.GetSchedule("")
		if err != nil {
			return scheduleLoadedMsg{err: err}
		}
		return scheduleLoadedMsg{schedule: resp}
	}
}

func (m Model) loadPlayerStats(playerID string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.GetPlayerStats(playerID)
		if err != nil {
			return playerStatsLoadedMsg{err: err}
		}
		return playerStatsLoadedMsg{stats: resp}
	}
}
