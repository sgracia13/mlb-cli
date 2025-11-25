package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"mlb-cli/internal/api"
)

var (
	// Flags for get subcommands
	seasonFlag string
	dateFlag   string
	teamFlag   string
)

// getCmd represents the get command group
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or more resources",
	Long: `Display one or more MLB resources.

Available resources:
  teams       List all MLB teams
  standings   Display division standings
  schedule    Show games for a specific date
  roster      Display a team's active roster

Examples:
  mlb get teams
  mlb get standings --season 2024
  mlb get schedule --date 2024-10-15
  mlb get roster --team LAD`,
	Aliases: []string{"g"},
}

// teamsCmd represents the 'get teams' command
var teamsCmd = &cobra.Command{
	Use:     "teams",
	Short:   "List all MLB teams",
	Long:    `Display a list of all MLB teams with their abbreviations and divisions.`,
	Aliases: []string{"team", "t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		teams, err := GetAPIClient().GetTeams()
		if err != nil {
			return fmt.Errorf("failed to get teams: %w", err)
		}
		return GetFormatter().PrintTeams(teams)
	},
}

// standingsCmd represents the 'get standings' command
var standingsCmd = &cobra.Command{
	Use:   "standings",
	Short: "Display division standings",
	Long: `Display MLB division standings for a specific season.

If no season is specified, the current year is used.

Examples:
  mlb get standings
  mlb get standings --season 2024
  mlb get standings -s 2023`,
	Aliases: []string{"standing", "stand", "st"},
	RunE: func(cmd *cobra.Command, args []string) error {
		season := seasonFlag
		if season == "" {
			season = time.Now().Format("2006")
		}

		standings, err := GetAPIClient().GetStandings(season)
		if err != nil {
			return fmt.Errorf("failed to get standings: %w", err)
		}
		return GetFormatter().PrintStandings(standings, season)
	},
}

// scheduleCmd represents the 'get schedule' command
var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Show games for a specific date",
	Long: `Display the MLB game schedule for a specific date.

If no date is specified, today's date is used.
Date format: YYYY-MM-DD

Examples:
  mlb get schedule
  mlb get schedule --date 2024-10-15
  mlb get schedule -d 2024-07-04`,
	Aliases: []string{"games", "sched", "sc"},
	RunE: func(cmd *cobra.Command, args []string) error {
		date := dateFlag
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		schedule, err := GetAPIClient().GetSchedule(date)
		if err != nil {
			return fmt.Errorf("failed to get schedule: %w", err)
		}
		return GetFormatter().PrintSchedule(schedule, date)
	},
}

// rosterCmd represents the 'get roster' command
var rosterCmd = &cobra.Command{
	Use:   "roster",
	Short: "Display a team's active roster",
	Long: `Display the active roster for an MLB team.

Specify the team using its abbreviation (e.g., LAD, NYY) or team ID.

Examples:
  mlb get roster --team LAD
  mlb get roster -t NYY
  mlb get roster --team 119

Team abbreviations:
  LAA, ARI, BAL, BOS, CHC, CIN, CLE, COL, DET, HOU,
  KC, LAD, WSH, NYM, OAK, PIT, SD, SEA, SF, STL,
  TB, TEX, TOR, MIN, PHI, ATL, CWS, MIA, NYY, MIL`,
	Aliases: []string{"rosters", "r"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if teamFlag == "" {
			return fmt.Errorf("team is required: use --team or -t flag")
		}

		teamID, err := api.ResolveTeamID(teamFlag)
		if err != nil {
			return err
		}

		roster, err := GetAPIClient().GetRoster(teamID)
		if err != nil {
			return fmt.Errorf("failed to get roster: %w", err)
		}
		return GetFormatter().PrintRoster(roster, teamID)
	},
}

func init() {
	// Add subcommands to 'get'
	getCmd.AddCommand(teamsCmd)
	getCmd.AddCommand(standingsCmd)
	getCmd.AddCommand(scheduleCmd)
	getCmd.AddCommand(rosterCmd)

	// Flags for standings
	standingsCmd.Flags().StringVarP(&seasonFlag, "season", "s", "",
		"Season year (default: current year)")

	// Flags for schedule
	scheduleCmd.Flags().StringVarP(&dateFlag, "date", "d", "",
		"Date in YYYY-MM-DD format (default: today)")

	// Flags for roster
	rosterCmd.Flags().StringVarP(&teamFlag, "team", "t", "",
		"Team abbreviation (e.g., LAD, NYY) or team ID")
	rosterCmd.MarkFlagRequired("team")
}
