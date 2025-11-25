package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"mlb-cli/internal/models"
)

const baseURL = "https://statsapi.mlb.com/api/v1"

// Client handles all MLB API interactions
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new MLB API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    baseURL,
	}
}

// fetch performs an HTTP GET request and returns the response body
func (c *Client) fetch(url string) ([]byte, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// GetTeams retrieves all MLB teams
func (c *Client) GetTeams() (*models.TeamsResponse, error) {
	data, err := c.fetch(c.baseURL + "/teams?sportId=1")
	if err != nil {
		return nil, err
	}

	var resp models.TeamsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetStandings retrieves standings for a given season
func (c *Client) GetStandings(season string) (*models.StandingsResponse, error) {
	url := fmt.Sprintf("%s/standings?leagueId=103,104&season=%s&standingsTypes=regularSeason", c.baseURL, season)
	data, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	var resp models.StandingsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetSchedule retrieves the game schedule for a given date
func (c *Client) GetSchedule(date string) (*models.ScheduleResponse, error) {
	url := fmt.Sprintf("%s/schedule?sportId=1&date=%s", c.baseURL, date)
	data, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	var resp models.ScheduleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// SearchPlayer searches for a player by name
func (c *Client) SearchPlayer(name string) (*models.PlayerSearchResponse, error) {
	url := fmt.Sprintf("%s/people/search?names=%s&sportId=1", c.baseURL, strings.ReplaceAll(name, " ", "%20"))
	data, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	var resp models.PlayerSearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetPlayerStats retrieves statistics for a player by ID
func (c *Client) GetPlayerStats(playerID string) (*models.PlayerStatsResponse, error) {
	statTypes := "yearByYear,career"
	url := fmt.Sprintf("%s/people/%s?hydrate=stats(group=[hitting,pitching],type=[%s])", c.baseURL, playerID, statTypes)
	data, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	var resp models.PlayerStatsResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetRoster retrieves the active roster for a team
func (c *Client) GetRoster(teamID string) (*models.RosterResponse, error) {
	url := fmt.Sprintf("%s/teams/%s/roster?rosterType=active", c.baseURL, teamID)
	data, err := c.fetch(url)
	if err != nil {
		return nil, err
	}

	var resp models.RosterResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// ResolveTeamID resolves a team abbreviation or ID to a team ID
func ResolveTeamID(input string) (string, error) {
	// First try to parse as int (direct ID)
	if _, err := strconv.Atoi(input); err == nil {
		return input, nil
	}

	// Try to look up by abbreviation
	upper := strings.ToUpper(input)
	if id, ok := models.GetTeamID(upper); ok {
		return strconv.Itoa(id), nil
	}

	return "", fmt.Errorf("unknown team: %s (use team abbreviation like LAD, NYY, or team ID)", input)
}
