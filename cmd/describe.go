package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Flags for describe subcommands
	statSeasonFlag string
)

// describeCmd represents the describe command group
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Show detailed information about a resource",
	Long: `Show detailed information about MLB resources.

Available resources:
  player   Search for and display player information
  stats    Display detailed statistics for a player

Examples:
  mlb describe player "Shohei Ohtani"
  mlb describe stats 660271
  mlb describe stats 660271 --season 2024`,
	Aliases: []string{"desc", "d"},
}

// playerCmd represents the 'describe player' command
var playerCmd = &cobra.Command{
	Use:   "player [name]",
	Short: "Search for and display player information",
	Long: `Search for an MLB player by name and display their information.

The search is case-insensitive and supports partial matches.

Examples:
  mlb describe player "Shohei Ohtani"
  mlb describe player ohtani
  mlb describe player "Mike Trout"`,
	Aliases: []string{"players", "p"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")

		players, err := GetAPIClient().SearchPlayer(name)
		if err != nil {
			return fmt.Errorf("failed to search player: %w", err)
		}
		return GetFormatter().PrintPlayer(players, name)
	},
}

// statsCmd represents the 'describe stats' command
var statsCmd = &cobra.Command{
	Use:   "stats [player_id]",
	Short: "Display detailed statistics for a player",
	Long: `Display detailed career and season statistics for an MLB player.

Use the player ID obtained from 'mlb describe player' command.
Optionally filter by season using the --season flag.

Examples:
  mlb describe stats 660271           # All career stats for Ohtani
  mlb describe stats 660271 --season 2024  # Only 2024 season
  mlb describe stats 545361 -s 2023   # Mike Trout's 2023 season`,
	Aliases: []string{"stat", "s"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		playerID := args[0]

		stats, err := GetAPIClient().GetPlayerStats(playerID)
		if err != nil {
			return fmt.Errorf("failed to get stats: %w", err)
		}
		return GetFormatter().PrintStats(stats, statSeasonFlag)
	},
}

func init() {
	// Add subcommands to 'describe'
	describeCmd.AddCommand(playerCmd)
	describeCmd.AddCommand(statsCmd)

	// Flags for stats
	statsCmd.Flags().StringVarP(&statSeasonFlag, "season", "s", "",
		"Filter stats by season year")
}
