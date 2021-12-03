package main

import (
	"net/http"

	"github.com/KunalDuran/duranzapi/module/cms"
	"github.com/KunalDuran/duranzapi/module/sports"
	"github.com/julienschmidt/httprouter"
)

// addRouteHandlers adds routes for various APIs.
func addRouteHandlers(router *httprouter.Router) {

	// General links
	router.GET("/", index)
	router.GET("/scorecard/:file", sports.GetScoreCard)

	// CMS
	router.GET("/cms/display", cms.DisplayPageCMS)

	router.GET("/cms/missing/players", cms.GetMissingPlayerDetails)
	router.GET("/cms/missing/teams", cms.GetMissingTeamDetails)
	router.GET("/cms/missing/venues", cms.GetMissingVenueDetails)

	router.POST("/cms/mapping/players", cms.MapPlayerDetails)
	router.POST("/cms/mapping/teams", cms.MapTeamDetails)
	router.POST("/cms/mapping/venues", cms.MapVenueDetails)

}

// Plain index page content
func index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Duranz API"))
}
