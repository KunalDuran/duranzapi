package data

import (
	"database/sql"

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
