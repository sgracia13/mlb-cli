package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the current view
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content string

	// Build the view
	content = m.renderHeader()
	content += m.renderTabs()
	content += m.renderBreadcrumb()

	if m.loading {
		content += m.renderLoading()
	} else if m.err != nil {
		content += m.renderError()
	} else {
		content += m.renderContent()
	}

	content += m.renderFilter()
	content += m.renderHelp()

	return AppStyle.Render(content)
}

func (m Model) renderHeader() string {
	title := "⚾ MLB CLI"
	return TitleStyle.Width(m.width - 4).Render(title) + "\n"
}

func (m Model) renderTabs() string {
	tabs := []string{"Teams", "Standings", "Schedule"}
	renderedTabs := make([]string, len(tabs))

	for i, tab := range tabs {
		if Tab(i) == m.currentTab {
			renderedTabs[i] = ActiveTabStyle.Render(tab)
		} else {
			renderedTabs[i] = InactiveTabStyle.Render(tab)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...) + "\n\n"
}

func (m Model) renderBreadcrumb() string {
	parts := []string{}

	switch m.currentView {
	case ViewTeams:
		parts = append(parts, "Teams")
	case ViewRoster:
		parts = append(parts, "Teams")
		if m.selectedTeam != nil {
			parts = append(parts, m.selectedTeam.Name)
		}
		parts = append(parts, "Roster")
	case ViewPlayer:
		parts = append(parts, "Teams")
		if m.selectedTeam != nil {
			parts = append(parts, m.selectedTeam.Name)
		}
		parts = append(parts, "Roster")
		if m.selectedPlayer != nil {
			parts = append(parts, m.selectedPlayer.Person.FullName)
		}
	case ViewStandings:
		parts = append(parts, "Standings")
	case ViewSchedule:
		parts = append(parts, "Schedule")
	}

	breadcrumb := strings.Join(parts, " > ")
	return BreadcrumbStyle.Render(breadcrumb) + "\n"
}

func (m Model) renderLoading() string {
	return fmt.Sprintf("\n%s Loading...\n", m.spinner.View())
}

func (m Model) renderError() string {
	return ErrorStyle.Render(fmt.Sprintf("\n✗ Error: %v\n\nPress 'r' to retry\n", m.err))
}

func (m Model) renderContent() string {
	switch m.currentView {
	case ViewTeams:
		return m.renderTeams()
	case ViewRoster:
		return m.renderRoster()
	case ViewPlayer:
		return m.renderPlayer()
	case ViewStandings:
		return m.renderStandings()
	case ViewSchedule:
		return m.renderSchedule()
	}
	return ""
}

func (m Model) renderTeams() string {
	if len(m.teams) == 0 {
		return MutedStyle.Render("No teams found")
	}

	var sb strings.Builder

	// Header
	header := fmt.Sprintf("%-5s %-25s %-25s", "ABBR", "TEAM", "DIVISION")
	sb.WriteString(TableHeaderStyle.Render(header) + "\n")

	// Calculate visible items based on height
	visibleItems := m.height - 12 // Account for header, tabs, breadcrumb, help, etc.
	if visibleItems < 5 {
		visibleItems = 5
	}

	// Get items to display (filtered or all)
	items := m.getTeamIndicesToShow()
	startIdx, endIdx := m.calculateVisibleRange(len(items), visibleItems)

	for i := startIdx; i < endIdx; i++ {
		teamIdx := items[i]
		if teamIdx >= len(m.teams) {
			continue
		}
		team := m.teams[teamIdx]
		line := fmt.Sprintf("%-5s %-25s %-25s", team.Abbreviation, team.Name, team.Division.Name)

		if i == m.cursor {
			sb.WriteString(SelectedStyle.Render(line) + "\n")
		} else {
			sb.WriteString(NormalStyle.Render(line) + "\n")
		}
	}

	// Show scroll indicator
	if len(items) > visibleItems {
		sb.WriteString(MutedStyle.Render(fmt.Sprintf("\n  Showing %d-%d of %d teams", startIdx+1, endIdx, len(items))))
	}

	return sb.String()
}

func (m Model) renderRoster() string {
	if len(m.roster) == 0 {
		return MutedStyle.Render("No players found")
	}

	var sb strings.Builder

	// Header
	header := fmt.Sprintf("%-4s %-25s %-6s %-15s", "#", "NAME", "POS", "STATUS")
	sb.WriteString(TableHeaderStyle.Render(header) + "\n")

	visibleItems := m.height - 12
	if visibleItems < 5 {
		visibleItems = 5
	}

	items := m.getRosterIndicesToShow()
	startIdx, endIdx := m.calculateVisibleRange(len(items), visibleItems)

	for i := startIdx; i < endIdx; i++ {
		rosterIdx := items[i]
		if rosterIdx >= len(m.roster) {
			continue
		}
		r := m.roster[rosterIdx]
		line := fmt.Sprintf("%-4s %-25s %-6s %-15s",
			r.JerseyNumber, r.Person.FullName, r.Position.Abbreviation, r.Status.Description)

		if i == m.cursor {
			sb.WriteString(SelectedStyle.Render(line) + "\n")
		} else {
			sb.WriteString(NormalStyle.Render(line) + "\n")
		}
	}

	if len(items) > visibleItems {
		sb.WriteString(MutedStyle.Render(fmt.Sprintf("\n  Showing %d-%d of %d players", startIdx+1, endIdx, len(items))))
	}

	return sb.String()
}

func (m Model) renderPlayer() string {
	if m.selectedPlayer == nil {
		return MutedStyle.Render("No player selected")
	}

	var sb strings.Builder
	p := m.selectedPlayer

	sb.WriteString(HeaderStyle.Render(p.Person.FullName) + "\n\n")
	sb.WriteString(fmt.Sprintf("  ID:       %d\n", p.Person.ID))
	sb.WriteString(fmt.Sprintf("  Number:   #%s\n", p.JerseyNumber))
	sb.WriteString(fmt.Sprintf("  Position: %s\n", p.Position.Abbreviation))
	sb.WriteString(fmt.Sprintf("  Status:   %s\n", p.Status.Description))

	// Show stats if loaded
	if m.playerStats != nil && len(m.playerStats.People) > 0 {
		player := m.playerStats.People[0]
		sb.WriteString("\n")

		// Consolidate stats by group type (hitting/pitching)
		type consolidatedStats struct {
			recentSeason string
			recentStat   map[string]interface{}
			careerStat   map[string]interface{}
		}
		statsByGroup := make(map[string]*consolidatedStats)

		for _, statGroup := range player.Stats {
			if len(statGroup.Splits) == 0 {
				continue
			}

			groupName := statGroup.Group.DisplayName
			if statsByGroup[groupName] == nil {
				statsByGroup[groupName] = &consolidatedStats{}
			}

			for i := range statGroup.Splits {
				split := &statGroup.Splits[i]
				if split.Season == "" {
					// Career stats
					statsByGroup[groupName].careerStat = split.Stat
				} else {
					// Season stats - keep the most recent
					if statsByGroup[groupName].recentSeason == "" || split.Season > statsByGroup[groupName].recentSeason {
						statsByGroup[groupName].recentSeason = split.Season
						statsByGroup[groupName].recentStat = split.Stat
					}
				}
			}
		}

		// Display consolidated stats in order: hitting first, then pitching
		for _, groupName := range []string{"hitting", "pitching"} {
			stats, ok := statsByGroup[groupName]
			if !ok || (stats.recentStat == nil && stats.careerStat == nil) {
				continue
			}

			sb.WriteString(HeaderStyle.Render(fmt.Sprintf("%s Stats", groupName)) + "\n")

			if stats.recentStat != nil {
				sb.WriteString(fmt.Sprintf("\n  %s Season:\n", stats.recentSeason))
				sb.WriteString(m.formatStats(groupName, stats.recentStat))
			}

			if stats.careerStat != nil {
				sb.WriteString("\n  Career:\n")
				sb.WriteString(m.formatStats(groupName, stats.careerStat))
			}

			sb.WriteString("\n")
		}
	}

	sb.WriteString(MutedStyle.Render("Press backspace to go back"))

	return sb.String()
}

func (m Model) formatStats(groupName string, stat map[string]interface{}) string {
	var sb strings.Builder

	if groupName == "hitting" {
		sb.WriteString(fmt.Sprintf("    AVG: %-8v  HR: %-6v  RBI: %-6v\n",
			getStatValue(stat, "avg"), getStatValue(stat, "homeRuns"), getStatValue(stat, "rbi")))
		sb.WriteString(fmt.Sprintf("    H:   %-8v  AB: %-6v  R:   %-6v\n",
			getStatValue(stat, "hits"), getStatValue(stat, "atBats"), getStatValue(stat, "runs")))
		sb.WriteString(fmt.Sprintf("    OBP: %-8v  SLG: %-6v  OPS: %-6v\n",
			getStatValue(stat, "obp"), getStatValue(stat, "slg"), getStatValue(stat, "ops")))
		sb.WriteString(fmt.Sprintf("    SB:  %-8v  BB: %-6v  SO:  %-6v\n",
			getStatValue(stat, "stolenBases"), getStatValue(stat, "baseOnBalls"), getStatValue(stat, "strikeOuts")))
	} else if groupName == "pitching" {
		sb.WriteString(fmt.Sprintf("    ERA:  %-8v  W: %-6v  L:   %-6v\n",
			getStatValue(stat, "era"), getStatValue(stat, "wins"), getStatValue(stat, "losses")))
		sb.WriteString(fmt.Sprintf("    G:    %-8v  GS: %-6v  SV:  %-6v\n",
			getStatValue(stat, "gamesPlayed"), getStatValue(stat, "gamesStarted"), getStatValue(stat, "saves")))
		sb.WriteString(fmt.Sprintf("    IP:   %-8v  SO: %-6v  BB:  %-6v\n",
			getStatValue(stat, "inningsPitched"), getStatValue(stat, "strikeOuts"), getStatValue(stat, "baseOnBalls")))
		sb.WriteString(fmt.Sprintf("    WHIP: %-8v  K/9: %-6v  BB/9: %-6v\n",
			getStatValue(stat, "whip"), getStatValue(stat, "strikeoutsPer9Inn"), getStatValue(stat, "walksPer9Inn")))
	}

	return sb.String()
}

func getStatValue(stat map[string]interface{}, key string) string {
	if val, ok := stat[key]; ok {
		return fmt.Sprintf("%v", val)
	}
	return "-"
}

func (m Model) renderStandings() string {
	if m.standings == nil || len(m.standings.Records) == 0 {
		return MutedStyle.Render("No standings data")
	}

	var sb strings.Builder

	for _, record := range m.standings.Records {
		sb.WriteString(HeaderStyle.Render(record.Division.Name) + "\n")

		header := fmt.Sprintf("%-3s %-22s %4s %4s %7s %6s", "#", "TEAM", "W", "L", "PCT", "GB")
		sb.WriteString(TableHeaderStyle.Render(header) + "\n")

		for _, tr := range record.TeamRecords {
			gb := tr.GamesBack
			if gb == "-" {
				gb = "—"
			}
			line := fmt.Sprintf("%-3s %-22s %4d %4d %7s %6s",
				tr.DivisionRank, tr.Team.Name, tr.Wins, tr.Losses, tr.WinningPct, gb)
			sb.WriteString(NormalStyle.Render(line) + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m Model) renderSchedule() string {
	if m.schedule == nil || len(m.schedule.Dates) == 0 {
		return MutedStyle.Render("No games scheduled")
	}

	var sb strings.Builder

	for _, date := range m.schedule.Dates {
		sb.WriteString(HeaderStyle.Render(fmt.Sprintf("📅 %s", date.Date)) + "\n")

		for _, game := range date.Games {
			status := game.Status.DetailedState
			matchup := fmt.Sprintf("%s @ %s", game.Teams.Away.Team.Name, game.Teams.Home.Team.Name)

			var line string
			if status == "Final" || status == "Game Over" {
				line = fmt.Sprintf("%-40s  %d - %d  [%s]",
					matchup, game.Teams.Away.Score, game.Teams.Home.Score, status)
			} else {
				line = fmt.Sprintf("%-40s  [%s]", matchup, status)
			}
			sb.WriteString(NormalStyle.Render(line) + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m Model) renderFilter() string {
	if !m.filterMode && m.filterText == "" {
		return ""
	}

	if m.filterMode {
		return "\n" + FilterStyle.Render("Filter: "+m.textInput.View())
	}

	return "\n" + FilterStyle.Render(fmt.Sprintf("Filter: %s (press / to change, esc to clear)", m.filterText))
}

func (m Model) renderHelp() string {
	var help string

	switch m.currentView {
	case ViewTeams:
		help = "↑/↓: Navigate  Enter: View Roster  Tab: Switch View  /: Filter  r: Refresh  q: Quit"
	case ViewRoster:
		help = "↑/↓: Navigate  Enter: Player Stats  Backspace: Back  /: Filter  r: Refresh  q: Quit"
	case ViewPlayer:
		help = "Backspace: Back  r: Refresh  q: Quit"
	case ViewStandings, ViewSchedule:
		help = "Tab: Switch View  r: Refresh  q: Quit"
	}

	return "\n" + HelpStyle.Render(help)
}

// Helper functions

func (m Model) getTeamIndicesToShow() []int {
	if len(m.filteredList) > 0 {
		return m.filteredList
	}
	indices := make([]int, len(m.teams))
	for i := range m.teams {
		indices[i] = i
	}
	return indices
}

func (m Model) getRosterIndicesToShow() []int {
	if len(m.filteredList) > 0 {
		return m.filteredList
	}
	indices := make([]int, len(m.roster))
	for i := range m.roster {
		indices[i] = i
	}
	return indices
}

func (m Model) calculateVisibleRange(total, visible int) (start, end int) {
	if total <= visible {
		return 0, total
	}

	// Keep cursor in the middle of visible range when possible
	half := visible / 2
	start = m.cursor - half
	if start < 0 {
		start = 0
	}

	end = start + visible
	if end > total {
		end = total
		start = end - visible
		if start < 0 {
			start = 0
		}
	}

	return start, end
}
