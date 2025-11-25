package models

// TeamsResponse represents the API response for teams endpoint
type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

// Team represents an MLB team
type Team struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Division     struct {
		Name string `json:"name"`
	} `json:"division"`
	League struct {
		Name string `json:"name"`
	} `json:"league"`
	Venue struct {
		Name string `json:"name"`
	} `json:"venue"`
}

// StandingsResponse represents the API response for standings endpoint
type StandingsResponse struct {
	Records []StandingsRecord `json:"records"`
}

// StandingsRecord represents a division's standings
type StandingsRecord struct {
	Division struct {
		Name string `json:"name"`
	} `json:"division"`
	TeamRecords []TeamRecord `json:"teamRecords"`
}

// TeamRecord represents a team's record in the standings
type TeamRecord struct {
	Team struct {
		Name string `json:"name"`
	} `json:"team"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	WinningPct   string `json:"winningPercentage"`
	GamesBack    string `json:"gamesBack"`
	DivisionRank string `json:"divisionRank"`
	Streak       struct {
		StreakCode string `json:"streakCode"`
	} `json:"streak"`
}

// ScheduleResponse represents the API response for schedule endpoint
type ScheduleResponse struct {
	Dates []ScheduleDate `json:"dates"`
}

// ScheduleDate represents a single date in the schedule
type ScheduleDate struct {
	Date  string         `json:"date"`
	Games []ScheduleGame `json:"games"`
}

// ScheduleGame represents a single game
type ScheduleGame struct {
	GamePk   int    `json:"gamePk"`
	GameDate string `json:"gameDate"`
	Status   struct {
		DetailedState string `json:"detailedState"`
	} `json:"status"`
	Teams struct {
		Away struct {
			Team struct {
				Name string `json:"name"`
			} `json:"team"`
			Score int `json:"score"`
		} `json:"away"`
		Home struct {
			Team struct {
				Name string `json:"name"`
			} `json:"team"`
			Score int `json:"score"`
		} `json:"home"`
	} `json:"teams"`
	Venue struct {
		Name string `json:"name"`
	} `json:"venue"`
}

// PlayerSearchResponse represents the API response for player search
type PlayerSearchResponse struct {
	People []Player `json:"people"`
}

// Player represents an MLB player
type Player struct {
	ID              int    `json:"id"`
	FullName        string `json:"fullName"`
	PrimaryPosition struct {
		Abbreviation string `json:"abbreviation"`
	} `json:"primaryPosition"`
	CurrentTeam struct {
		Name string `json:"name"`
	} `json:"currentTeam"`
	BatSide struct {
		Code string `json:"code"`
	} `json:"batSide"`
	PitchHand struct {
		Code string `json:"code"`
	} `json:"pitchHand"`
	BirthDate string `json:"birthDate"`
	Height    string `json:"height"`
	Weight    int    `json:"weight"`
	Active    bool   `json:"active"`
}

// PlayerStatsResponse represents the API response for player stats
type PlayerStatsResponse struct {
	People []PlayerWithStats `json:"people"`
}

// PlayerWithStats represents a player with their statistics
type PlayerWithStats struct {
	FullName string      `json:"fullName"`
	Stats    []StatGroup `json:"stats"`
}

// StatGroup represents a group of statistics (hitting/pitching)
type StatGroup struct {
	Group struct {
		DisplayName string `json:"displayName"`
	} `json:"group"`
	Splits []StatSplit `json:"splits"`
}

// StatSplit represents statistics for a specific season
type StatSplit struct {
	Season string                 `json:"season"`
	Stat   map[string]interface{} `json:"stat"`
}

// RosterResponse represents the API response for team roster
type RosterResponse struct {
	Roster []RosterEntry `json:"roster"`
}

// RosterEntry represents a single roster entry
type RosterEntry struct {
	Person struct {
		ID       int    `json:"id"`
		FullName string `json:"fullName"`
	} `json:"person"`
	Position struct {
		Abbreviation string `json:"abbreviation"`
	} `json:"position"`
	JerseyNumber string `json:"jerseyNumber"`
	Status       struct {
		Description string `json:"description"`
	} `json:"status"`
}

// TeamAbbreviations maps team abbreviations to their IDs
var TeamAbbreviations = map[string]int{
	"LAA": 108, "ARI": 109, "BAL": 110, "BOS": 111, "CHC": 112,
	"CIN": 113, "CLE": 114, "COL": 115, "DET": 116, "HOU": 117,
	"KC":  118, "LAD": 119, "WSH": 120, "NYM": 121, "OAK": 133,
	"PIT": 134, "SD":  135, "SEA": 136, "SF":  137, "STL": 138,
	"TB":  139, "TEX": 140, "TOR": 141, "MIN": 142, "PHI": 143,
	"ATL": 144, "CWS": 145, "MIA": 146, "NYY": 147, "MIL": 158,
}

// GetTeamID returns the team ID for a given abbreviation or ID string
func GetTeamID(input string) (int, bool) {
	if id, ok := TeamAbbreviations[input]; ok {
		return id, true
	}
	return 0, false
}
