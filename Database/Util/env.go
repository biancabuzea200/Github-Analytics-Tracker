package Util

import (
	"GitHub-Analytics-Tracker/Util"
)

var ()

// GetDSN retrieves the database (postgres DSN) environment variable
func GetDSN() (string, error) {
	return Util.GetENV("POSTGRES_DSN")
}

// GetSchema retrieves the database schema (directory) where all DB operations should happen in
func GetSchema() (string, error) {
	return Util.GetENV("DB_SCHEMA")
}
