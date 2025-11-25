package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"mlb-cli/internal/models"
)

// Format represents the output format type
type Format string

const (
	FormatTable Format = "table"
	FormatWide  Format = "wide"
	FormatJSON  Format = "json"
)

// ParseFormat parses a string into a Format type
func ParseFormat(s string) Format {
	switch strings.ToLower(s) {
	case "json":
		return FormatJSON
	case "wide":
		return FormatWide
	default:
		return FormatTable
	}
}

// Formatter handles output formatting
type Formatter struct {
	format Format
}

// NewFormatter creates a new formatter with the specified format
func NewFormatter(format Format) *Formatter {
	return &Formatter{format: format}
}

// PrintTeams outputs teams in the specified format
func (f *Formatter) PrintTeams(teams *models.TeamsResponse) error {
	if f.format == FormatJSON {
		return printJSON(teams)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	fmt.Println()
	fmt.Println("⚾ MLB Teams")
	fmt.Println(strings.Repeat("─", 70))

	if f.format == FormatWide {
		fmt.Fprintf(w, "ID\tABBR\tNAME\tDIVISION\tLEAGUE\tVENUE\n")
		fmt.Fprintln(w, strings.Repeat("─", 70))
		for _, t := range teams.Teams {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
				t.ID, t.Abbreviation, t.Name, t.Division.Name, t.League.Name, t.Venue.Name)
		}
	} else {
		fmt.Fprintf(w, "ABBR\tNAME\tDIVISION\n")
		fmt.Fprintln(w, strings.Repeat("─", 70))
		for _, t := range teams.Teams {
			fmt.Fprintf(w, "%s\t%s\t%s\n", t.Abbreviation, t.Name, t.Division.Name)
		}
	}

	return nil
}

// PrintStandings outputs standings in the specified format
func (f *Formatter) PrintStandings(standings *models.StandingsResponse, season string) error {
	if f.format == FormatJSON {
		return printJSON(standings)
	}

	fmt.Printf("\n⚾ MLB Standings - %s\n", season)

	for _, rec := range standings.Records {
		fmt.Printf("\n%s\n", rec.Division.Name)
		fmt.Println(strings.Repeat("─", 75))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		if f.format == FormatWide {
			fmt.Fprintf(w, "#\tTEAM\tW\tL\tPCT\tGB\tSTREAK\n")
		} else {
			fmt.Fprintf(w, "#\tTEAM\tW\tL\tPCT\tGB\tSTREAK\n")
		}

		for _, tr := range rec.TeamRecords {
			gb := tr.GamesBack
			if gb == "-" {
				gb = "—"
			}
			fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%s\t%s\t%s\n",
				tr.DivisionRank,
				tr.Team.Name,
				tr.Wins,
				tr.Losses,
				tr.WinningPct,
				gb,
				tr.Streak.StreakCode,
			)
		}
		w.Flush()
	}

	return nil
}

// PrintSchedule outputs schedule in the specified format
func (f *Formatter) PrintSchedule(schedule *models.ScheduleResponse, date string) error {
	if f.format == FormatJSON {
		return printJSON(schedule)
	}

	fmt.Printf("\n⚾ MLB Games - %s\n", date)
	fmt.Println(strings.Repeat("─", 80))

	if len(schedule.Dates) == 0 {
		fmt.Println("No games scheduled for this date.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	if f.format == FormatWide {
		fmt.Fprintf(w, "GAME ID\tMATCHUP\tSCORE\tSTATUS\tVENUE\n")
		fmt.Fprintln(w, strings.Repeat("─", 80))
	}

	for _, d := range schedule.Dates {
		for _, g := range d.Games {
			status := g.Status.DetailedState
			matchup := fmt.Sprintf("%s @ %s", g.Teams.Away.Team.Name, g.Teams.Home.Team.Name)

			if f.format == FormatWide {
				score := "-"
				if status == "Final" || status == "Game Over" {
					score = fmt.Sprintf("%d - %d", g.Teams.Away.Score, g.Teams.Home.Score)
				}
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", g.GamePk, matchup, score, status, g.Venue.Name)
			} else {
				if status == "Final" || status == "Game Over" {
					fmt.Fprintf(w, "%s\t%d - %d\t[%s]\n",
						matchup, g.Teams.Away.Score, g.Teams.Home.Score, status)
				} else {
					fmt.Fprintf(w, "%s\t\t[%s]\n", matchup, status)
				}
			}
		}
	}

	return nil
}

// PrintPlayer outputs player info in the specified format
func (f *Formatter) PrintPlayer(players *models.PlayerSearchResponse, searchName string) error {
	if f.format == FormatJSON {
		return printJSON(players)
	}

	fmt.Printf("\n⚾ Player Search: \"%s\"\n", searchName)
	fmt.Println(strings.Repeat("─", 70))

	if len(players.People) == 0 {
		fmt.Println("No players found.")
		return nil
	}

	for _, p := range players.People {
		status := "Active"
		if !p.Active {
			status = "Inactive"
		}
		team := p.CurrentTeam.Name
		if team == "" {
			team = "Free Agent"
		}

		fmt.Printf("\n%s (ID: %d)\n", p.FullName, p.ID)
		fmt.Printf("  Position: %s | Team: %s | Status: %s\n",
			p.PrimaryPosition.Abbreviation, team, status)
		fmt.Printf("  Bats: %s | Throws: %s | Height: %s | Weight: %d lbs\n",
			p.BatSide.Code, p.PitchHand.Code, p.Height, p.Weight)

		if f.format == FormatWide && p.BirthDate != "" {
			fmt.Printf("  Born: %s\n", p.BirthDate)
		}
	}

	if len(players.People) > 0 {
		fmt.Printf("\n💡 Tip: Use 'mlb describe stats %d' to see career stats\n", players.People[0].ID)
	}

	return nil
}

// PrintStats outputs player stats in the specified format
func (f *Formatter) PrintStats(stats *models.PlayerStatsResponse, season string) error {
	if f.format == FormatJSON {
		return printJSON(stats)
	}

	if len(stats.People) == 0 {
		fmt.Println("Player not found.")
		return nil
	}

	player := stats.People[0]
	fmt.Printf("\n⚾ Stats for %s\n", player.FullName)

	for _, statGroup := range player.Stats {
		if len(statGroup.Splits) == 0 {
			continue
		}

		groupName := statGroup.Group.DisplayName
		fmt.Printf("\n%s Stats:\n", groupName)
		fmt.Println(strings.Repeat("─", 80))

		for _, split := range statGroup.Splits {
			if season != "" && split.Season != season && split.Season != "" {
				continue
			}

			seasonLabel := split.Season
			if seasonLabel == "" {
				seasonLabel = "Career"
			}

			fmt.Printf("\n  %s:\n", seasonLabel)

			stat := split.Stat
			if groupName == "hitting" {
				printStatLine("    AVG", stat, "avg")
				printStatLine("    HR", stat, "homeRuns")
				printStatLine("    RBI", stat, "rbi")
				printStatLine("    H", stat, "hits")
				printStatLine("    AB", stat, "atBats")
				printStatLine("    OBP", stat, "obp")
				printStatLine("    SLG", stat, "slg")
				printStatLine("    OPS", stat, "ops")
				printStatLine("    SB", stat, "stolenBases")
				printStatLine("    BB", stat, "baseOnBalls")
				printStatLine("    SO", stat, "strikeOuts")
			} else if groupName == "pitching" {
				printStatLine("    ERA", stat, "era")
				printStatLine("    W", stat, "wins")
				printStatLine("    L", stat, "losses")
				printStatLine("    G", stat, "gamesPlayed")
				printStatLine("    GS", stat, "gamesStarted")
				printStatLine("    SV", stat, "saves")
				printStatLine("    IP", stat, "inningsPitched")
				printStatLine("    SO", stat, "strikeOuts")
				printStatLine("    BB", stat, "baseOnBalls")
				printStatLine("    WHIP", stat, "whip")
			}
		}
	}

	return nil
}

// PrintRoster outputs team roster in the specified format
func (f *Formatter) PrintRoster(roster *models.RosterResponse, teamID string) error {
	if f.format == FormatJSON {
		return printJSON(roster)
	}

	fmt.Printf("\n⚾ Active Roster (Team ID: %s)\n", teamID)
	fmt.Println(strings.Repeat("─", 65))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	if f.format == FormatWide {
		fmt.Fprintf(w, "#\tID\tNAME\tPOS\tSTATUS\n")
	} else {
		fmt.Fprintf(w, "#\tNAME\tPOS\tSTATUS\n")
	}
	fmt.Fprintln(w, strings.Repeat("─", 65))

	for _, r := range roster.Roster {
		if f.format == FormatWide {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\n",
				r.JerseyNumber, r.Person.ID, r.Person.FullName,
				r.Position.Abbreviation, r.Status.Description)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				r.JerseyNumber, r.Person.FullName,
				r.Position.Abbreviation, r.Status.Description)
		}
	}

	return nil
}

// printJSON outputs data as formatted JSON
func printJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// printStatLine prints a single stat line
func printStatLine(label string, stat map[string]interface{}, key string) {
	if val, ok := stat[key]; ok {
		fmt.Printf("%-8s %v\n", label+":", val)
	}
}
