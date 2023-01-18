package sports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/KunalDuran/duranzapi/module/data"
	"github.com/KunalDuran/duranzapi/module/util"

	"github.com/julienschmidt/httprouter"
)

var DATASET_BASE = "/home/kunalduran/Desktop/duranz_api/dev/all_json/"

// GetCricsheetData : Reads the match json file
func GetCricsheetData(f_path string) (data.Match, error) {
	var matchData data.Match
	body, err := ioutil.ReadFile(f_path)
	if err != nil {
		return matchData, err
	}

	err = json.Unmarshal(body, &matchData)
	if err != nil {
		return matchData, err
	}
	return matchData, nil
}

func GetScoreCard(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	jsonID := p.ByName("file")
	match, err := GetCricsheetData(DATASET_BASE + jsonID + `.json`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var objScoreCard data.ScoreCard

	var AllInnings []data.Innings
	for _, inning := range match.Innings {
		fmt.Println("Scorecard process started innings for : ", inning.Team)

		var objInning data.Innings
		objInning.InningDetail = inning.Team

		var objExtra data.Extras
		var objBatsman = map[string]data.Batting{}
		var objBowler = map[string]data.Bowling{}
		var batsmanCount, bowlerCount, runningScore, wicketCnt int
		var fowArr []string
		var overRuns = map[int]int{}
		for _, over := range inning.Overs {
			overRuns[over.Over] = 0
			for _, delivery := range over.Deliveries {

				// Score Calculations
				runningScore += delivery.Runs.Total
				overRuns[over.Over] += delivery.Runs.Total

				// Batsman Init
				if _, exist := objBatsman[delivery.Batter]; !exist {
					batsmanCount++
					var tempBat data.Batting
					tempBat.ID = batsmanCount
					tempBat.Name = delivery.Batter
					objBatsman[delivery.Batter] = tempBat
				}
				batsman := objBatsman[delivery.Batter]
				batsman.Runs += delivery.Runs.Batter
				batsman.Balls += 1

				// Bowler Init
				if _, exist := objBowler[delivery.Bowler]; !exist {
					bowlerCount++
					var tempBowler data.Bowling
					tempBowler.ID = bowlerCount
					tempBowler.Name = delivery.Bowler
					objBowler[delivery.Bowler] = tempBowler
				}
				bowler := objBowler[delivery.Bowler]
				bowler.Runs += delivery.Runs.Batter
				bowler.Balls += 1

				if delivery.Runs.Batter == 4 {
					batsman.Fours++
				} else if delivery.Runs.Batter == 6 {
					batsman.Sixes++
				}

				// Calculate Extras
				if delivery.Extras != (data.Extras{}) {
					if delivery.Extras.Byes > 0 {
						objExtra.Byes += delivery.Extras.Byes
						overRuns[over.Over] -= delivery.Extras.Byes
					} else if delivery.Extras.LegByes > 0 {
						objExtra.LegByes += delivery.Extras.LegByes
						overRuns[over.Over] -= delivery.Extras.LegByes
					} else if delivery.Extras.NoBall > 0 {
						// remove ball count if No Ball
						batsman.Balls -= 1
						bowler.Balls -= 1
						bowler.Runs += delivery.Extras.NoBall
						objExtra.NoBall += delivery.Extras.NoBall
					} else if delivery.Extras.Wides > 0 {
						// remove ball count if Wide Ball
						batsman.Balls -= 1
						bowler.Balls -= 1
						bowler.Runs += delivery.Extras.Wides
						objExtra.Wides += delivery.Extras.Wides
					}
				}

				// Check for Wicket
				for _, wicket := range delivery.Wickets {
					if wicket.Kind != "" && wicket.PlayerOut != "" {
						batsman.Out = wicket.Kind
						wicketCnt++
						fowStr := fmt.Sprint(wicketCnt, "-", runningScore, "(", wicket.PlayerOut, ")")
						fowArr = append(fowArr, fowStr)

						// bowler
						if wicket.Kind != "run out" {
							bowler.Wickets++
						}
					}
				}

				// bind all info and calculations
				objBatsman[delivery.Batter] = batsman
				objBowler[delivery.Bowler] = bowler
			}

			// check maiden over
			if val, ok := overRuns[over.Over]; ok && val == 0 {
				if len(over.Deliveries) > 0 {
					bowler := objBowler[over.Deliveries[0].Bowler]
					bowler.Maiden++
					objBowler[over.Deliveries[0].Bowler] = bowler
				}
			}
		}

		var allBatsman []data.Batting
		for _, batter := range objBatsman {
			if batter.Balls > 0 {
				batter.StrikeRate = util.Round((float64(batter.Runs)*100)/float64(batter.Balls), 0.01, 2)
			}
			if batter.Out == "" {
				batter.Out = "not out"
			}
			allBatsman = append(allBatsman, batter)
		}

		var allBowler []data.Bowling
		for _, bowler := range objBowler {
			if bowler.Balls > 0 {
				bowler.Economy = util.Round(float64(bowler.Runs)/(float64(bowler.Balls)/float64(6)), 0.01, 2)
			}
			bowler.Overs = fmt.Sprint(bowler.Balls/6) + "." + fmt.Sprint(bowler.Balls%6)
			allBowler = append(allBowler, bowler)
		}

		objExtra.Total = objExtra.Byes + objExtra.LegByes + objExtra.Wides + objExtra.NoBall
		objInning.Extras = objExtra
		objInning.Batting = allBatsman
		objInning.FallOfWickets = strings.Join(fowArr, " , ")
		objInning.Bowling = allBowler

		AllInnings = append(AllInnings, objInning)
	}

	resultStr := match.Info.Outcome.Winner + " Won by "
	if match.Info.Outcome.By.Runs > 0 {
		resultStr += strconv.Itoa(match.Info.Outcome.By.Runs) + " Runs"
	} else if match.Info.Outcome.By.Wickets > 0 {
		resultStr += strconv.Itoa(match.Info.Outcome.By.Wickets) + " Wickets"
	}
	objScoreCard.Result = resultStr
	objScoreCard.Innings = AllInnings

	final := util.JSONMessageWrappedObj(http.StatusOK, objScoreCard)
	util.WebResponseJSONObject(w, r, http.StatusOK, final)

	// ioutil.WriteFile(`C:\Users\Kunal\Desktop\Duranz\duranz_api\scoreCard.json`, strScoreCard, 0777)
}

func PlayerStatsAPI(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	playerName := p.ByName("player")
	bio := r.URL.Query().Get("bio")
	format := util.CleanText(r.URL.Query().Get("format"), true)
	season := util.CleanText(r.URL.Query().Get("season"), true)
	vsteam := util.CleanText(r.URL.Query().Get("vsteam"), true)

	if playerName == "all" {
		// do something
	}
	if bio != "true" {
		bio = ""
	}

	vsTeamID := 0
	if vsteam != "" {
		vsTeamID = data.GetTeamID(vsteam)
	}
	var playerFinalAll []data.PlayerStatsExt
	objAllPlayerStats := data.GetPlayerStats(playerName, format, season, vsTeamID)

	for pname, pstats := range objAllPlayerStats {
		var playerFinal data.PlayerStatsExt

		// General Stats
		playerFinal.TeamID = pstats[0].TeamID.Int64
		playerFinal.PlayerID = pstats[0].PlayerID.Int64
		playerFinal.PlayerName = pname

		// playerFinal.InningsID += pstat.InningsID.Int64
		// playerFinal.OutBowler += pstat.OutBowler.Int64
		// playerFinal.OutFielder += pstat.OutFielder.Int64
		// playerFinal.OutType += pstat.OutType.Int64
		playerFinal.SeasonID = season
		// playerFinal.SeasonType += pstat.SeasonType.Int64

		for _, pstat := range pstats {

			// bind batting stats
			playerFinal.Batting.BallsFaced += pstat.BallsFaced.Int64
			// playerFinal.Batting.BattingOrder += pstat.BattingOrder.Int64
			playerFinal.Batting.DotBallsPlayed += pstat.DotBallsPlayed.Int64
			playerFinal.Batting.Doubles += pstat.Doubles.Int64
			playerFinal.Batting.FoursHit += pstat.FoursHit.Int64
			playerFinal.Batting.RunsScored += pstat.RunsScored.Int64
			playerFinal.Batting.Singles += pstat.Singles.Int64
			playerFinal.Batting.SixesHit += pstat.SixesHit.Int64
			playerFinal.Batting.Triples += pstat.Triples.Int64
			playerFinal.Batting.IsBatted += pstat.IsBatted.Int64
			if pstat.RunsScored.Int64 >= 100 {
				playerFinal.Batting.Hundreds++
			} else if pstat.RunsScored.Int64 >= 50 {
				playerFinal.Batting.Fifties++
			}
			if pstat.OutType.String == "not out" {
				playerFinal.Batting.NotOuts++
			}
			if pstat.RunsScored.Int64 == 0 && pstat.OutType.String != "not out" && pstat.OutType.String != "" {
				playerFinal.Batting.Ducks++
			}

			if pstat.RunsScored.Int64 > playerFinal.Batting.HighestScore {
				playerFinal.Batting.HighestScore = pstat.RunsScored.Int64
			}

			// bind bowling stats
			// playerFinal.Bowling.BowlingOrder += pstat.BowlingOrder.Int64
			playerFinal.Bowling.DotsBowled += pstat.DotsBowled.Int64
			playerFinal.Bowling.MaidenOver += pstat.MaidenOver.Int64
			playerFinal.Bowling.BallsBowled += pstat.BallsBowled.Int64
			playerFinal.Bowling.ExtrasConceded += pstat.ExtrasConceded.Int64
			playerFinal.Bowling.FoursConceded += pstat.FoursConceded.Int64
			playerFinal.Bowling.RunsConceded += pstat.RunsConceded.Int64
			playerFinal.Bowling.SixesConceded += pstat.SixesConceded.Int64
			playerFinal.Bowling.WicketsTaken += pstat.WicketsTaken.Int64
			if pstat.WicketsTaken.Int64 >= 5 {
				playerFinal.Bowling.Fifers++
			}
			if playerFinal.Bowling.BestBowling != "" {
				bowlingFigures := strings.Split(playerFinal.Bowling.BestBowling, "/")
				wickets, runs := bowlingFigures[0], bowlingFigures[1]
				wicketsInt, _ := strconv.ParseInt(wickets, 10, 64)
				runsInt, _ := strconv.ParseInt(runs, 10, 64)
				if wicketsInt < pstat.WicketsTaken.Int64 {
					playerFinal.Bowling.BestBowling = fmt.Sprint(pstat.WicketsTaken.Int64) + "/" + fmt.Sprint(pstat.RunsConceded.Int64)
				} else if wicketsInt == pstat.WicketsTaken.Int64 && runsInt > pstat.RunsConceded.Int64 {
					playerFinal.Bowling.BestBowling = fmt.Sprint(pstat.WicketsTaken.Int64) + "/" + fmt.Sprint(pstat.RunsConceded.Int64)
				}
			} else {
				playerFinal.Bowling.BestBowling = fmt.Sprint(pstat.WicketsTaken.Int64) + "/" + fmt.Sprint(pstat.RunsConceded.Int64)
			}

			// bind fieling stats
			playerFinal.Fielding.Catches += pstat.Catches.Int64
			playerFinal.Fielding.Stumpings += pstat.Stumpings.Int64
			playerFinal.Fielding.RunOut += pstat.RunOut.Int64
		}
		if playerFinal.Bowling.BallsBowled > 0 {
			playerFinal.Bowling.OversBowled = fmt.Sprint(playerFinal.Bowling.BallsBowled/6) + "." + fmt.Sprint(playerFinal.Bowling.BallsBowled%6)
		}

		if playerFinal.Batting.IsBatted-playerFinal.Batting.NotOuts > 0 {
			playerFinal.Batting.Average = util.Round((float64(playerFinal.Batting.RunsScored))/float64(playerFinal.Batting.IsBatted-playerFinal.Batting.NotOuts), 0.01, 2)
		}

		if playerFinal.Batting.BallsFaced > 0 {
			playerFinal.Batting.StrikeRate = util.Round((float64(playerFinal.Batting.RunsScored)*100)/float64(playerFinal.Batting.BallsFaced), 0.01, 2)
		}

		if playerFinal.Bowling.WicketsTaken > 0 {
			playerFinal.Bowling.Average = util.Round((float64(playerFinal.Bowling.RunsConceded))/float64(playerFinal.Bowling.WicketsTaken), 0.01, 2)
		}

		if playerFinal.Bowling.BallsBowled > 0 {
			playerFinal.Bowling.Economy = util.Round((float64(playerFinal.Bowling.RunsConceded))/(float64(playerFinal.Bowling.BallsBowled)/6), 0.01, 2)
		}

		playerFinalAll = append(playerFinalAll, playerFinal)
	}

	final := util.JSONMessageWrappedObj(http.StatusOK, playerFinalAll)
	util.WebResponseJSONObject(w, r, http.StatusOK, final)

}

func TeamStatsAPI(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	teamName := p.ByName("team")
	gender := util.CleanText(r.URL.Query().Get("gender"), true)
	if gender == "" {
		gender = "male"
	}
	season := util.CleanText(r.URL.Query().Get("season"), true)

	teamID := data.GetTeamID(teamName)
	objAllTeamStats := data.GetTeamStats(teamID, gender, season)
	var teamFinalAll data.DuranzTeamStats

	for _, objTeam := range objAllTeamStats {

		if objTeam.TossWinner.Valid {
			if objTeam.TossWinner.Int64 == int64(teamID) {
				teamFinalAll.TossWin++
			}
		}

		if objTeam.WinningTeam.Valid {
			if objTeam.WinningTeam.Int64 == int64(teamID) {
				teamFinalAll.MatchWin++

				// count if team won while batting first or chasing first
				if objTeam.TossWinner.Valid {
					if objTeam.TossWinner.Int64 == int64(teamID) { // team won the toss
						if objTeam.TossDecision.String == "bat" {
							teamFinalAll.BatFirstWin++
						} else {
							teamFinalAll.ChasingWin++
						}
					} else {
						if objTeam.TossDecision.String == "bat" { // other team won the toss
							teamFinalAll.ChasingWin++
						} else {
							teamFinalAll.BatFirstWin++
						}
					}
				}
			}
		}
	}

	teamFinalAll.TotalMatches = len(objAllTeamStats)
	teamFinalAll.MatchWinPercent = util.Round(float64(teamFinalAll.MatchWin)/float64(teamFinalAll.TotalMatches), 0.01, 2)
	teamFinalAll.TossWinPercent = util.Round(float64(teamFinalAll.TossWin)/float64(teamFinalAll.TotalMatches), 0.01, 2)
	teamFinalAll.ChasingWinPer = math.Round((float64(teamFinalAll.ChasingWin)/float64(teamFinalAll.MatchWin))*10000) / 100
	teamFinalAll.BatFirstWinPer = math.Round((float64(teamFinalAll.BatFirstWin)/float64(teamFinalAll.MatchWin))*10000) / 100
	//AvgScore /Inn
	//Highest Score
	//Lowest Score

	final := util.JSONMessageWrappedObj(http.StatusOK, teamFinalAll)
	util.WebResponseJSONObject(w, r, http.StatusOK, final)

}

func BatsmanVSBowlerAPI(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//
}

func PlayerList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	playerList := data.GetPlayerList()
	final := util.JSONMessageWrappedObj(http.StatusOK, playerList)
	util.WebResponseJSONObject(w, r, http.StatusOK, final)

}

func TeamList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	teamList := data.GetTeamList()
	final := util.JSONMessageWrappedObj(http.StatusOK, teamList)
	util.WebResponseJSONObject(w, r, http.StatusOK, final)
}
