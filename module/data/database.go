package data

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/KunalDuran/duranzapi/module/sports"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// SportsDb is the pointer to the duranz database resource.
var SportsDb *sql.DB

// InitDB initialises the database pools with
func InitDB(host, port, user, password string) (sportsDb *sql.DB, err error) {
	SportsDb, err = sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/duranz")

	if err != nil {
		return nil, err
	}
	if err = SportsDb.Ping(); err != nil {
		return nil, err
	}
	return SportsDb, nil
}

// RequestAPIData - Calls a (API) URL and return the data from the request.
func RequestAPIData(url string, headers map[string]string) ([]byte, int, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 500, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := http.Client{}
	client.Timeout = time.Duration(120 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	statusCode := resp.StatusCode

	// read body
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, statusCode, err
	}
	return body, statusCode, nil
}

func GetPlayerDetails() []sports.PlayerDetailsInt {
	var objAllPlayer = []sports.PlayerDetailsInt{}

	sqlStr := `SELECT * FROM duranz_cricket_players`

	rows, err := SportsDb.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var objPlayer sports.PlayerDetailsInt
		err = rows.Scan(
			&objPlayer.PlayerID,
			&objPlayer.PlayerName,
			&objPlayer.DisplayName,
			&objPlayer.FirstName,
			&objPlayer.LastName,
			&objPlayer.ShortName,
			&objPlayer.UniqueShortName,
			&objPlayer.DOB,
			&objPlayer.BattingStyle,
			&objPlayer.BowlingStyle,
			&objPlayer.IsOverseas,
			&objPlayer.CricSheetID,
			&objPlayer.DateAdded,
			&objPlayer.Status,
		)
		if err != nil {
			panic(err)
		}
		objAllPlayer = append(objAllPlayer, objPlayer)
	}
	return objAllPlayer
}

func GetTeamDetails() []sports.TeamDetailsInt {
	var objAllTeams = []sports.TeamDetailsInt{}

	sqlStr := `SELECT * FROM duranz_teams`

	rows, err := SportsDb.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var objTeam sports.TeamDetailsInt
		err = rows.Scan(
			&objTeam.TeamID,
			&objTeam.TeamName,
			&objTeam.TeamType,
			&objTeam.FilterName,
			&objTeam.ABBR,
			&objTeam.TeamColor,
			&objTeam.Icon,
			&objTeam.URL,
			&objTeam.Jersey,
			&objTeam.Flag,
			&objTeam.Status,
			&objTeam.DateAdded,
		)
		if err != nil {
			panic(err)
		}
		objAllTeams = append(objAllTeams, objTeam)
	}
	return objAllTeams
}

func GetVenueDetails() []sports.VenueDetailsInt {
	var objAllVenue = []sports.VenueDetailsInt{}

	sqlStr := `SELECT * FROM duranz_venue`

	rows, err := SportsDb.Query(sqlStr)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var objVenue sports.VenueDetailsInt
		err = rows.Scan(
			&objVenue.VenueID,
			&objVenue.Venue,
			&objVenue.FilterName,
			&objVenue.FriendlyName,
			&objVenue.City,
			&objVenue.Country,
			&objVenue.State,
			&objVenue.StateABBR,
			&objVenue.OfficialTeam,
			&objVenue.Capacity,
			&objVenue.Dimension,
			&objVenue.Opened,
			&objVenue.Description,
			&objVenue.ShortName,
			&objVenue.TimeZone,
			&objVenue.Weather,
			&objVenue.PitchType,
			&objVenue.DateAdded,
			&objVenue.Status,
		)
		if err != nil {
			panic(err)
		}
		objAllVenue = append(objAllVenue, objVenue)
	}
	return objAllVenue
}

func InsertPlayerDetails(objPlayer sports.PlayerDetailsExt) {

	sqlStr := `INSERT INTO duranz_cricket_players(player_id,player_name,display_name,first_name,last_name,
		short_name,unique_short_name,dob,batting_style_1_id,bowling_style_1_id,is_overseas,cricsheet_id,
		date_added,status)
		VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		player_id=VALUES(player_id),player_name=VALUES(player_name),display_name=VALUES(display_name),
		first_name=VALUES(first_name),last_name=VALUES(last_name),short_name=VALUES(short_name),
		unique_short_name=VALUES(unique_short_name),dob=VALUES(dob),batting_style_1_id=VALUES(batting_style_1_id),
		bowling_style_1_id=VALUES(bowling_style_1_id),is_overseas=VALUES(is_overseas),
		cricsheet_id=VALUES(cricsheet_id),date_added=VALUES(date_added),status=VALUES(status)`

	_, err := SportsDb.Exec(
		sqlStr,
		objPlayer.PlayerID,
		objPlayer.PlayerName,
		objPlayer.DisplayName,
		objPlayer.FirstName,
		objPlayer.LastName,
		objPlayer.ShortName,
		objPlayer.UniqueShortName,
		objPlayer.DOB,
		objPlayer.BattingStyle,
		objPlayer.BowlingStyle,
		objPlayer.IsOverseas,
		objPlayer.CricSheetID,
		objPlayer.DateAdded,
		objPlayer.Status,
	)
	if err != nil {
		panic(err)
	}
}

func InsertTeamDetails(objTeam sports.TeamDetailsExt) {

	sqlStr := `INSERT INTO duranz_teams(
		team_id,team_name,team_type,filtername,abbreviation,team_color,icon,url,jersey,flag,status,dateadded) 
		VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) 
		ON DUPLICATE KEY UPDATE 
		team_id=VALUES(team_id),team_name=VALUES(team_name),team_type=VALUES(team_type),
		filtername=VALUES(filtername),abbreviation=VALUES(abbreviation),team_color=VALUES(team_color),
		icon=VALUES(icon),url=VALUES(url),jersey=VALUES(jersey),flag=VALUES(flag),
		status=VALUES(status),dateadded=VALUES(dateadded)`

	_, err := SportsDb.Exec(
		sqlStr,
		objTeam.TeamID,
		objTeam.TeamName,
		objTeam.TeamType,
		objTeam.FilterName,
		objTeam.ABBR,
		objTeam.TeamColor,
		objTeam.Icon,
		objTeam.URL,
		objTeam.Jersey,
		objTeam.Flag,
		objTeam.Status,
		objTeam.DateAdded,
	)
	if err != nil {
		panic(err)
	}
}

func InsertVenueDetails(objVenue sports.VenueDetailsExt) {

	sqlStr := `INSERT INTO duranz_venue (
		venue_id,venue,filtername,friendlyname,city,country,state,state_abbr,official_team,capacity,
		dimensions,opened,description,shortname,timezone,weather,pitch_type,dateadded,status)
	 	VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 	
		venue_id=VALUES(venue_id),venue=VALUES(venue),filtername=VALUES(filtername),
		friendlyname=VALUES(friendlyname),city=VALUES(city),country=VALUES(country),
		state=VALUES(state),state_abbr=VALUES(state_abbr),official_team=VALUES(official_team),
		capacity=VALUES(capacity),dimensions=VALUES(dimensions),opened=VALUES(opened),
		description=VALUES(description),shortname=VALUES(shortname),timezone=VALUES(timezone),
		weather=VALUES(weather),pitch_type=VALUES(pitch_type),dateadded=VALUES(dateadded),
		status=VALUES(status)
	`

	_, err := SportsDb.Exec(
		sqlStr,
		objVenue.VenueID,
		objVenue.Venue,
		objVenue.FilterName,
		objVenue.FriendlyName,
		objVenue.City,
		objVenue.Country,
		objVenue.State,
		objVenue.StateABBR,
		objVenue.OfficialTeam,
		objVenue.Capacity,
		objVenue.Dimension,
		objVenue.Opened,
		objVenue.Description,
		objVenue.ShortName,
		objVenue.TimeZone,
		objVenue.Weather,
		objVenue.PitchType,
		objVenue.DateAdded,
		objVenue.Status,
	)
	if err != nil {
		panic(err)
	}
}
