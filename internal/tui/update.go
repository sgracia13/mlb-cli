package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle filter mode input
		if m.filterMode {
			return m.handleFilterInput(msg)
		}

		// Handle normal key input
		return m.handleKeyInput(msg)

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case teamsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.teams = msg.teams
			m.resetFilter()
		}

	case rosterLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.roster = msg.roster
			m.resetFilter()
		}

	case standingsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.standings = msg.standings
		}

	case scheduleLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.schedule = msg.schedule
		}

	case playerStatsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.playerStats = msg.stats
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(msg, m.keys.Quit):
		return m, tea.Quit

	case matchKey(msg, m.keys.Up):
		m.moveCursor(-1)

	case matchKey(msg, m.keys.Down):
		m.moveCursor(1)

	case matchKey(msg, m.keys.Enter):
		return m.handleSelect()

	case matchKey(msg, m.keys.Back):
		return m.handleBack()

	case matchKey(msg, m.keys.Filter):
		m.filterMode = true
		m.textInput.Focus()
		return m, nil

	case matchKey(msg, m.keys.Tab):
		return m.nextTab()

	case matchKey(msg, m.keys.ShiftTab):
		return m.prevTab()

	case matchKey(msg, m.keys.Refresh):
		return m.refresh()
	}

	return m, nil
}

func (m Model) handleFilterInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter, tea.KeyEscape:
		m.filterMode = false
		m.textInput.Blur()
		if msg.Type == tea.KeyEscape {
			m.textInput.SetValue("")
			m.resetFilter()
		} else {
			m.applyFilter()
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	m.filterText = m.textInput.Value()
	m.applyFilter()
	return m, cmd
}

func (m *Model) moveCursor(delta int) {
	listLen := m.currentListLength()
	if listLen == 0 {
		return
	}

	m.cursor += delta
	if m.cursor < 0 {
		m.cursor = listLen - 1
	} else if m.cursor >= listLen {
		m.cursor = 0
	}
}

func (m Model) currentListLength() int {
	if len(m.filteredList) > 0 {
		return len(m.filteredList)
	}

	switch m.currentView {
	case ViewTeams:
		return len(m.teams)
	case ViewRoster:
		return len(m.roster)
	case ViewStandings:
		if m.standings == nil {
			return 0
		}
		count := 0
		for _, r := range m.standings.Records {
			count += len(r.TeamRecords)
		}
		return count
	case ViewSchedule:
		if m.schedule == nil {
			return 0
		}
		count := 0
		for _, d := range m.schedule.Dates {
			count += len(d.Games)
		}
		return count
	}
	return 0
}

func (m Model) handleSelect() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case ViewTeams:
		if len(m.teams) > 0 {
			idx := m.getActualIndex()
			if idx < len(m.teams) {
				m.selectedTeam = &m.teams[idx]
				m.pushHistory()
				m.currentView = ViewRoster
				m.cursor = 0
				m.loading = true
				m.resetFilter()
				return m, tea.Batch(
					m.spinner.Tick,
					m.loadRoster(fmt.Sprintf("%d", m.selectedTeam.ID)),
				)
			}
		}

	case ViewRoster:
		if len(m.roster) > 0 {
			idx := m.getActualIndex()
			if idx < len(m.roster) {
				m.selectedPlayer = &m.roster[idx]
				m.playerStats = nil // Clear previous stats
				m.pushHistory()
				m.currentView = ViewPlayer
				m.cursor = 0
				m.loading = true
				return m, tea.Batch(
					m.spinner.Tick,
					m.loadPlayerStats(fmt.Sprintf("%d", m.selectedPlayer.Person.ID)),
				)
			}
		}
	}

	return m, nil
}

func (m Model) handleBack() (tea.Model, tea.Cmd) {
	if len(m.history) > 0 {
		m.currentView = m.history[len(m.history)-1]
		m.history = m.history[:len(m.history)-1]
		m.cursor = 0
		m.resetFilter()
	}
	return m, nil
}

func (m Model) nextTab() (tea.Model, tea.Cmd) {
	m.currentTab = Tab((int(m.currentTab) + 1) % 3)
	return m.switchToTab()
}

func (m Model) prevTab() (tea.Model, tea.Cmd) {
	m.currentTab = Tab((int(m.currentTab) + 2) % 3)
	return m.switchToTab()
}

func (m Model) switchToTab() (tea.Model, tea.Cmd) {
	m.cursor = 0
	m.history = make([]View, 0)
	m.resetFilter()

	switch m.currentTab {
	case TabTeams:
		m.currentView = ViewTeams
		if len(m.teams) == 0 {
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, m.loadTeams)
		}
	case TabStandings:
		m.currentView = ViewStandings
		if m.standings == nil {
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, m.loadStandings())
		}
	case TabSchedule:
		m.currentView = ViewSchedule
		if m.schedule == nil {
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, m.loadSchedule())
		}
	}

	return m, nil
}

func (m Model) refresh() (tea.Model, tea.Cmd) {
	m.loading = true
	m.err = nil

	switch m.currentView {
	case ViewTeams:
		return m, tea.Batch(m.spinner.Tick, m.loadTeams)
	case ViewRoster:
		if m.selectedTeam != nil {
			return m, tea.Batch(m.spinner.Tick, m.loadRoster(fmt.Sprintf("%d", m.selectedTeam.ID)))
		}
	case ViewPlayer:
		if m.selectedPlayer != nil {
			return m, tea.Batch(m.spinner.Tick, m.loadPlayerStats(fmt.Sprintf("%d", m.selectedPlayer.Person.ID)))
		}
	case ViewStandings:
		return m, tea.Batch(m.spinner.Tick, m.loadStandings())
	case ViewSchedule:
		return m, tea.Batch(m.spinner.Tick, m.loadSchedule())
	}

	return m, nil
}

func (m *Model) pushHistory() {
	m.history = append(m.history, m.currentView)
}

func (m *Model) resetFilter() {
	m.filterText = ""
	m.filteredList = nil
	m.textInput.SetValue("")
}

func (m *Model) applyFilter() {
	if m.filterText == "" {
		m.filteredList = nil
		return
	}

	filter := strings.ToLower(m.filterText)
	m.filteredList = make([]int, 0)

	switch m.currentView {
	case ViewTeams:
		for i, t := range m.teams {
			if strings.Contains(strings.ToLower(t.Name), filter) ||
				strings.Contains(strings.ToLower(t.Abbreviation), filter) {
				m.filteredList = append(m.filteredList, i)
			}
		}
	case ViewRoster:
		for i, r := range m.roster {
			if strings.Contains(strings.ToLower(r.Person.FullName), filter) ||
				strings.Contains(strings.ToLower(r.Position.Abbreviation), filter) {
				m.filteredList = append(m.filteredList, i)
			}
		}
	}

	// Reset cursor if out of bounds
	if m.cursor >= len(m.filteredList) && len(m.filteredList) > 0 {
		m.cursor = 0
	}
}

func (m Model) getActualIndex() int {
	if len(m.filteredList) > 0 && m.cursor < len(m.filteredList) {
		return m.filteredList[m.cursor]
	}
	return m.cursor
}

func matchKey(msg tea.KeyMsg, binding interface{}) bool {
	switch b := binding.(type) {
	case key.Binding:
		return key.Matches(msg, b)
	case string:
		return msg.String() == b
	}
	return false
}
