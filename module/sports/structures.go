package sports

import "database/sql"

type Match struct {
	Meta struct {
		DataVersion string `json:"data_version"`
		Created     string `json:"created"`
		Revision    int    `json:"revision"`
	} `json:"meta"`
	Info struct {
		BallsPerOver    int      `json:"balls_per_over"`
		City            string   `json:"city"`
		Dates           []string `json:"dates"`
		Gender          string   `json:"gender"`
		MatchType       string   `json:"match_type"`
		MatchTypeNumber int      `json:"match_type_number"`
		Officials       struct {
			MatchReferees  []string `json:"match_referees"`
			ReserveUmpires []string `json:"reserve_umpires"`
			TvUmpires      []string `json:"tv_umpires"`
			Umpires        []string `json:"umpires"`
		} `json:"officials"`
		Outcome struct {
			By struct {
				Runs    int `json:"runs"`
				Wickets int `json:"wickets"`
			} `json:"by"`
			Winner string `json:"winner"`
		} `json:"outcome"`
		Overs         int      `json:"overs"`
		PlayerOfMatch []string `json:"player_of_match"`
		Players       struct {
			Pakistan   []string `json:"Pakistan"`
			WestIndies []string `json:"West Indies"`
		} `json:"players"`
		Register Registry `json:"registry"`
		Season   string   `json:"season"`
		TeamType string   `json:"team_type"`
		Teams    []string `json:"teams"`
		Toss     struct {
			Decision string `json:"decision"`
			Winner   string `json:"winner"`
		} `json:"toss"`
		Venue string `json:"venue"`
	} `json:"info"`
	Innings []struct {
		Team  string `json:"team"`
		Overs []struct {
			Over       int `json:"over"`
			Deliveries []struct {
				Batter     string `json:"batter"`
				Bowler     string `json:"bowler"`
				Extras     `json:"extras"`
				NonStriker string `json:"non_striker"`
				Runs       struct {
					Batter int `json:"batter"`
					Extras int `json:"extras"`
					Total  int `json:"total"`
				} `json:"runs"`
				Wickets []struct {
					Kind      string `json:"kind"`
					PlayerOut string `json:"player_out"`
				}
			} `json:"deliveries"`
		} `json:"overs"`
		Powerplays []struct {
			From float64 `json:"from"`
			To   float64 `json:"to"`
			Type string  `json:"type"`
		} `json:"powerplays"`
		Target struct {
			Overs int `json:"overs"`
			Runs  int `json:"runs"`
		} `json:"target,omitempty"`
	} `json:"innings"`
}

type Registry struct {
	People map[string]string `json:"people"`
}

// SCORE CARD STRUCTS
type ScoreCard struct {
	Innings []Innings `json:"innings"`
	Result  string    `json:"result"`
}

type Innings struct {
	InningID      int       `json:"innings_id"`
	InningDetail  string    `json:"innings_detail"`
	Bowling       []Bowling `json:"bowling"`
	Batting       []Batting `json:"batting"`
	Extras        `json:"extras"`
	FallOfWickets string `json:"fall_of_wickets"`
}

type Bowling struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Overs   string  `json:"overs"`
	Maiden  int     `json:"maiden"`
	Runs    int     `json:"runs"`
	Wickets int     `json:"wickets"`
	Economy float64 `json:"economy"`
	Balls   int     `json:"-"`
}

type Batting struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Runs       int     `json:"runs"`
	Balls      int     `json:"balls"`
	Fours      int     `json:"fours"`
	Sixes      int     `json:"sixes"`
	StrikeRate float64 `json:"strike_rate"`
	Out        string  `json:"out"`
}

type Extras struct {
	Wides   int `json:"wides"`
	NoBall  int `json:"noballs"`
	Byes    int `json:"byes"`
	LegByes int `json:"legbyes"`
	Total   int `json:"total"`
}

type PlayerDetailsInt struct {
	PlayerID        sql.NullInt64
	PlayerName      sql.NullString
	DisplayName     sql.NullString
	FirstName       sql.NullString
	LastName        sql.NullString
	ShortName       sql.NullString
	UniqueShortName sql.NullString
	DOB             sql.NullString
	BattingStyle    sql.NullInt64
	BowlingStyle    sql.NullInt64
	IsOverseas      sql.NullInt64
	CricSheetID     sql.NullString
	DateAdded       sql.NullString
	Status          sql.NullInt64
}

type PlayerDetailsExt struct {
	PlayerID        int64  `json:"player_id"`
	PlayerName      string `json:"player_name"`
	DisplayName     string `json:"display_name"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ShortName       string `json:"short_name"`
	UniqueShortName string `json:"unique_short_name"`
	DOB             string `json:"dob"`
	BattingStyle    int64  `json:"batting_style"`
	BowlingStyle    int64  `json:"bowling_style"`
	IsOverseas      int64  `json:"is_overseas"`
	CricSheetID     string `json:"cricsheet_id"`
	DateAdded       string `json:"date_added"`
	Status          int64  `json:"status"`
}

type TeamDetailsInt struct {
	TeamID     sql.NullInt64
	TeamName   sql.NullString
	TeamType   sql.NullString
	FilterName sql.NullString
	ABBR       sql.NullString
	TeamColor  sql.NullString
	Icon       sql.NullString
	URL        sql.NullString
	Jersey     sql.NullString
	Flag       sql.NullString
	Status     sql.NullInt64
	DateAdded  sql.NullString
}

type TeamDetailsExt struct {
	TeamID     int64  `json:"team_id"`
	TeamName   string `json:"team_name"`
	TeamType   string `json:"team_type"`
	FilterName string `json:"filter_name"`
	ABBR       string `json:"abbr"`
	TeamColor  string `json:"team_color"`
	Icon       string `json:"icon"`
	URL        string `json:"url"`
	Jersey     string `json:"jersey"`
	Flag       string `json:"flag"`
	Status     int64  `json:"status"`
	DateAdded  string `json:"date_added"`
}

type VenueDetailsInt struct {
	VenueID      sql.NullInt64
	Venue        sql.NullString
	FilterName   sql.NullString
	FriendlyName sql.NullString
	City         sql.NullString
	Country      sql.NullString
	State        sql.NullString
	StateABBR    sql.NullString
	OfficialTeam sql.NullString
	Capacity     sql.NullString
	Dimension    sql.NullString
	Opened       sql.NullString
	Description  sql.NullString
	ShortName    sql.NullString
	TimeZone     sql.NullString
	Weather      sql.NullString
	PitchType    sql.NullString
	DateAdded    sql.NullString
	Status       sql.NullInt64
}

type VenueDetailsExt struct {
	VenueID      int64  `json:"venue_id"`
	Venue        string `json:"venue_name"`
	FilterName   string `json:"filter_name"`
	FriendlyName string `json:"friendly_name"`
	City         string `json:"city"`
	Country      string `json:"country"`
	State        string `json:"state"`
	StateABBR    string `json:"state_abbr"`
	OfficialTeam string `json:"official_team"`
	Capacity     string `json:"capacity"`
	Dimension    string `json:"dimension"`
	Opened       string `json:"opened"`
	Description  string `json:"description"`
	ShortName    string `json:"short_name"`
	TimeZone     string `json:"time_zone"`
	Weather      string `json:"weather"`
	PitchType    string `json:"pitch_type"`
	DateAdded    string `json:"date_added"`
	Status       int64  `json:"status"`
}
