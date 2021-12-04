package cms

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/KunalDuran/duranzapi/module/data"
	"github.com/KunalDuran/duranzapi/module/sports"
	"github.com/KunalDuran/duranzapi/module/util"
	"github.com/julienschmidt/httprouter"
)

func DisplayPageCMS(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	param := r.URL.Query().Get("page")
	if param != "players" && param != "venues" && param != "teams" {
		w.Write([]byte("This is CMS Page, Pass players, venues or teams in query param 'page'"))
		return
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	body, _, err := data.RequestAPIData("http://localhost:5000/cms/missing/"+param, headers)
	if err != nil {
		panic(err)
	}

	tmpl, err := template.ParseFiles("templates/" + param + ".html")
	if err != nil {
		panic(err)
	}

	switch param {
	case "players":

		finalPlayers := struct {
			Content []sports.PlayerDetailsExt `json:"content"`
		}{}
		err = json.Unmarshal(body, &finalPlayers)
		fmt.Println(len(finalPlayers.Content))
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, finalPlayers.Content)

	case "teams":

		finalTeams := struct {
			Content []sports.PlayerDetailsExt `json:"content"`
		}{}
		err = json.Unmarshal(body, &finalTeams)
		fmt.Println(len(finalTeams.Content))
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, finalTeams.Content)

	case "venues":

		finalVenues := struct {
			Content []sports.VenueDetailsExt `json:"content"`
		}{}
		err = json.Unmarshal(body, &finalVenues)
		fmt.Println(len(finalVenues.Content))
		if err != nil {
			panic(err)
		}
		tmpl.Execute(w, finalVenues.Content[0])
	}
}

func GetMissingPlayerDetails(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objPlayers := data.GetPlayerDetails()

	var finalPlayers []sports.PlayerDetailsExt
	for _, objPlayer := range objPlayers {

		var finalPlayer sports.PlayerDetailsExt

		finalPlayer.PlayerID = objPlayer.PlayerID.Int64
		finalPlayer.PlayerName = objPlayer.PlayerName.String
		finalPlayer.DisplayName = objPlayer.DisplayName.String
		finalPlayer.FirstName = objPlayer.FirstName.String
		finalPlayer.LastName = objPlayer.LastName.String
		finalPlayer.ShortName = objPlayer.ShortName.String
		finalPlayer.UniqueShortName = objPlayer.UniqueShortName.String
		finalPlayer.DOB = objPlayer.DOB.String
		finalPlayer.BattingStyle = objPlayer.BattingStyle.String
		finalPlayer.BowlingStyle = objPlayer.BowlingStyle.String
		finalPlayer.IsOverseas = objPlayer.IsOverseas.Int64
		finalPlayer.CricSheetID = objPlayer.CricSheetID.String
		finalPlayer.DateAdded = objPlayer.DateAdded.String
		finalPlayer.Status = objPlayer.Status.Int64

		finalPlayers = append(finalPlayers, finalPlayer)
	}

	final := util.JSONMessageWrappedObj(http.StatusOK, finalPlayers)
	util.WebResponseJSONObject(w, r, http.StatusOK, final)
}

func GetMissingTeamDetails(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objTeams := data.GetTeamDetails()

	var finalTeams []sports.TeamDetailsExt
	for _, objTeam := range objTeams {

		var finalTeam sports.TeamDetailsExt

		finalTeam.TeamID = objTeam.TeamID.Int64
		finalTeam.TeamName = objTeam.TeamName.String
		finalTeam.TeamType = objTeam.TeamType.String
		finalTeam.FilterName = objTeam.FilterName.String
		finalTeam.ABBR = objTeam.ABBR.String
		finalTeam.TeamColor = objTeam.TeamColor.String
		finalTeam.Icon = objTeam.Icon.String
		finalTeam.URL = objTeam.URL.String
		finalTeam.Jersey = objTeam.Jersey.String
		finalTeam.Flag = objTeam.Flag.String
		finalTeam.Status = objTeam.Status.Int64
		finalTeam.DateAdded = objTeam.DateAdded.String

		finalTeams = append(finalTeams, finalTeam)
	}

	final := util.JSONMessageWrappedObj(http.StatusOK, finalTeams)
	util.WebResponseJSONObject(w, r, http.StatusOK, final)
}

func GetMissingVenueDetails(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	objVenues := data.GetVenueDetails()

	var finalVenues []sports.VenueDetailsExt
	for _, objVenue := range objVenues {

		var finalVenue sports.VenueDetailsExt
		finalVenue.VenueID = objVenue.VenueID.Int64
		finalVenue.Venue = objVenue.Venue.String
		finalVenue.FilterName = objVenue.FilterName.String
		finalVenue.FriendlyName = objVenue.FriendlyName.String
		finalVenue.City = objVenue.City.String
		finalVenue.Country = objVenue.Country.String
		finalVenue.State = objVenue.State.String
		finalVenue.StateABBR = objVenue.StateABBR.String
		finalVenue.OfficialTeam = objVenue.OfficialTeam.String
		finalVenue.Capacity = objVenue.Capacity.String
		finalVenue.Dimension = objVenue.Dimension.String
		finalVenue.Opened = objVenue.Opened.String
		finalVenue.Description = objVenue.Description.String
		finalVenue.ShortName = objVenue.ShortName.String
		finalVenue.TimeZone = objVenue.TimeZone.String
		finalVenue.Weather = objVenue.Weather.String
		finalVenue.PitchType = objVenue.PitchType.String
		finalVenue.DateAdded = objVenue.DateAdded.String
		finalVenue.Status = objVenue.Status.Int64
		finalVenues = append(finalVenues, finalVenue)
	}

	final := util.JSONMessageWrappedObj(http.StatusOK, finalVenues)
	util.WebResponseJSONObject(w, r, http.StatusOK, final)
}

func MapVenueDetails(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var finalVenue sports.VenueDetailsExt
	venueIDInt, _ := strconv.ParseInt(r.FormValue("venue_id"), 10, 64)
	finalVenue.VenueID = venueIDInt
	finalVenue.Venue = r.FormValue("venue_name")
	finalVenue.FilterName = r.FormValue("filter_name")
	finalVenue.FriendlyName = r.FormValue("friendly_name")
	finalVenue.City = r.FormValue("city")
	finalVenue.Country = r.FormValue("country")
	finalVenue.State = r.FormValue("state")
	finalVenue.StateABBR = r.FormValue("state_abbr")
	finalVenue.OfficialTeam = r.FormValue("official_team")
	finalVenue.Capacity = r.FormValue("capacity")
	finalVenue.Dimension = r.FormValue("dimension")
	finalVenue.Opened = r.FormValue("opened")
	finalVenue.Description = r.FormValue("description")
	finalVenue.ShortName = r.FormValue("shortname")
	finalVenue.TimeZone = r.FormValue("timezone")
	finalVenue.Weather = r.FormValue("weather")
	finalVenue.PitchType = r.FormValue("pitch_type")
	time.Now().Format("2006-01-02 15:04:05")
	finalVenue.DateAdded = time.Now().Format("2006-01-02 15:04:05")
	venueStatus, _ := strconv.ParseInt(r.FormValue("status"), 10, 64)
	finalVenue.Status = venueStatus

	data.InsertVenueDetails(finalVenue)
	http.Redirect(w, r, "/cms/missing/venues", 200)
}

func MapPlayerDetails(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var finalPlayer sports.PlayerDetailsExt

	playerIDInt, _ := strconv.ParseInt(r.FormValue(""), 10, 64)
	finalPlayer.PlayerID = playerIDInt
	finalPlayer.PlayerName = r.Form.Get("")
	finalPlayer.DisplayName = r.Form.Get("")
	finalPlayer.FirstName = r.Form.Get("")
	finalPlayer.LastName = r.Form.Get("")
	finalPlayer.ShortName = r.Form.Get("")
	finalPlayer.UniqueShortName = r.Form.Get("")
	finalPlayer.DOB = r.Form.Get("")
	finalPlayer.BattingStyle = r.Form.Get("")
	finalPlayer.BowlingStyle = r.Form.Get("")
	overseasInt, _ := strconv.ParseInt(r.FormValue(""), 10, 64)
	finalPlayer.IsOverseas = overseasInt
	finalPlayer.CricSheetID = r.Form.Get("")
	finalPlayer.DateAdded = r.Form.Get("")
	statusInt, _ := strconv.ParseInt(r.FormValue(""), 10, 64)
	finalPlayer.Status = statusInt

	data.InsertPlayerDetails(finalPlayer)
}

func MapTeamDetails(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var finalTeam sports.TeamDetailsExt

	teamIDInt, _ := strconv.ParseInt(r.FormValue(""), 10, 64)
	finalTeam.TeamID = teamIDInt
	finalTeam.TeamName = r.Form.Get("")
	finalTeam.TeamType = r.Form.Get("")
	finalTeam.FilterName = r.Form.Get("")
	finalTeam.ABBR = r.Form.Get("")
	finalTeam.TeamColor = r.Form.Get("")
	finalTeam.Icon = r.Form.Get("")
	finalTeam.URL = r.Form.Get("")
	finalTeam.Jersey = r.Form.Get("")
	finalTeam.Flag = r.Form.Get("")
	statusInt, _ := strconv.ParseInt(r.FormValue(""), 10, 64)
	finalTeam.Status = statusInt
	finalTeam.DateAdded = r.Form.Get("")

	data.InsertTeamDetails(finalTeam)
}
