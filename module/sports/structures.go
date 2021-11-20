package sports

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
