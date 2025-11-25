package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"mlb-cli/internal/api"
	"mlb-cli/internal/output"
	"mlb-cli/internal/tui"
)

var (
	outputFormat string

	apiClient *api.Client
	formatter *output.Formatter
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mlb",
	Short: "MLB Stats CLI - Your command-line interface for MLB data",
	Long: `MLB Stats CLI is a command-line tool for accessing MLB statistics and information.

Similar to kubectl for Kubernetes, mlb provides a familiar interface for
querying teams, players, standings, schedules, and statistics.

Running 'mlb' without arguments launches the interactive TUI (like k9s).

Examples:
  mlb                               # Launch interactive TUI
  mlb get teams                     # List all MLB teams
  mlb get standings --season 2024   # View standings for 2024
  mlb get schedule                  # Today's games
  mlb get roster --team LAD         # Dodgers roster

  mlb describe player "Shohei Ohtani"  # Search for a player
  mlb describe stats 660271            # Player stats by ID`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize shared instances before each command
		apiClient = api.NewClient()
		formatter = output.NewFormatter(output.ParseFormat(outputFormat))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Launch interactive TUI when no subcommand is provided
		if err := tui.Run(); err != nil {
			return fmt.Errorf("failed to start TUI: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table",
		"Output format: table, wide, or json")

	// Add command groups
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(versionCmd)
}

// GetAPIClient returns the shared API client
func GetAPIClient() *api.Client {
	return apiClient
}

// GetFormatter returns the shared formatter
func GetFormatter() *output.Formatter {
	return formatter
}
