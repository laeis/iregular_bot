package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
type DbIniter interface{
	InitDb()
}
type AppDb struct {
	db   *sql.DB
	path string
}

func (a *AppDb) GetDbPath() string {
	return a.path
}

func (a *AppDb) SetDbPath(path string) {
	a.path = path
}

func (a *AppDb) InitDb() error {
	db, err := sql.Open("sqlite3", a.path)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	a.db = db
	return nil
}

func (a *AppDb) CloseDb() (err error) {
	err = a.db.Close()
	return
}