package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/KunalDuran/duranzapi/module/data"
	_ "github.com/go-sql-driver/mysql"

	"github.com/julienschmidt/httprouter"
)

func main() {

	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "root"
	dbPass := "password"

	var env = strings.ToLower(os.Getenv("Environment"))
	if env == "production" {
		fmt.Println("working in Production")
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
	log.Fatal(http.ListenAndServe(":5000", router))
}
