package main

import (
	"irregular_bot/help"
	"os"
	"path/filepath"
	"strings"
)

const (
	dbDefaultFolder = "/storage/"
	dbDefaultName   = "verb.db"
)

func init() {
	var dataBaseName string
	if name := os.Getenv("DataBase"); name != "" {
		dataBaseName = name
		if !strings.Contains(".db", name) {
			dataBaseName += ".db"
		}

	} else {
		dataBaseName = dbDefaultName
	}
	appPath, err := filepath.Abs(".")
	help.CheckErr(err)

	databasePath := appPath + dbDefaultFolder

	DatabaseLink = databasePath + dataBaseName

	_, err = os.Stat(DatabaseLink)
	if os.IsNotExist(err) {
		f, err := os.Create(DatabaseLink)
		help.CheckErr(err)
		defer f.Close()
	}

}
