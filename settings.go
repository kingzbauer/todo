package main

import (
	"os"
)

type Settings struct {
	DBURI  string
	DBName string
}

var settings *Settings

func FromEnvironment() *Settings {
	if settings != nil {
		return settings
	}

	s = &Settings{}

	s.DBURI = os.Getenv("DB_URI")
	s.DBName = os.Getenv("DB_NAME")

	settings = s

	return settings
}
