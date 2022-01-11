package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/KunalDuran/duranzapi/module/data"
	"github.com/KunalDuran/duranzapi/module/sports"
	_ "github.com/go-sql-driver/mysql"

	"github.com/julienschmidt/httprouter"
)

func main() {

	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "root"
	dbPass := "password"
	port := ":5000"

	if runtime.GOOS == "windows" {
		dbPass = ""
		sports.DATASET_BASE = `C:\Users\Kunal\Desktop\Duranz\duranz_api\matchdata\odis_json\`
	}

	var env = strings.ToLower(os.Getenv("Environment"))
	if env == "production" {
		dbUser = "kunal"
		// port = ":80"
		fmt.Println("working in Production")
		sports.DATASET_BASE = `/home/ubuntu/duranz/matchdata/`
	} else if env == "development" {
		fmt.Println("working in Development")
	}

	// Connect the database
	_, err := data.InitDB(dbHost, dbPort, dbUser, dbPass)
	if err != nil {
		log.Panic(err)
	}

	router := httprouter.New()
	router.RedirectTrailingSlash = true
	addRouteHandlers(router)

	fmt.Println("Duranz API initialized")
	log.Fatal(http.ListenAndServe(port, router))
}
