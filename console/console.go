package main

// use for parse and import verbs from site to app table
import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
)

const Url = "http://englishstyle.net/grammar/verb/irregular-verbs/"

func main() {

	db, err := sql.Open("sqlite3", "/home/art/go/src/irregular_bot/storage/verb.db")
	if err != nil {
		log.Fatal("Error database conection. ", err)
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 'verbs' (
		'id' INTEGER PRIMARY KEY AUTOINCREMENT,
		'present' varchar(32) ,
		'p_transcription'  varchar(32),
		'past_simple'  varchar(32),
		'ps_transcription'  varchar(32),
		'past_partisiple'  varchar(32),
		'pp_transription' varchar(32),
		'created_at' TIMESTAMP
		DEFAULT CURRENT_TIMESTAMP,
		'updated_at' TIMESTAMP
		DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Error database create table. ", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS 'translate' (
		'id' INTEGER PRIMARY KEY AUTOINCREMENT,
		'verb_id' INTEGER NOT NULL,
		'ru_Ru' varchar(64) ,
		'created_at' TIMESTAMP
		DEFAULT CURRENT_TIMESTAMP,
		'updated_at' TIMESTAMP
		DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (verb_id) 
		REFERENCES verbs(id)
		ON UPDATE CASCADE 
		ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal("Error database create table. ", err)
	}
	response, err := http.Get(Url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	document.Find(".entry table tr").Each(func(i int, tr *goquery.Selection) {
		list := make([]interface{}, 0, 6)
		tr.Find("td").Each(func(i int, td *goquery.Selection) {
			list = append(list, strings.TrimSpace(td.Text()))
		})
		if len(list) == 0 {
			return
		}

		stmt, err := db.Prepare(`
				INSERT INTO verbs(
					present,
					p_transcription,
					past_simple,
					ps_transcription,
					past_partisiple,
					pp_transription
					) values(?,?,?,?,?,?)
			`)
		if err != nil {
			log.Fatal("Error prepare query in db. ", err)
		}
		res, err := stmt.Exec(list[:6]...)
		if err != nil {
			log.Fatal("Error insert in verbs db. ", err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Fatal("Error get last insert id. ", err)
		}

		stmt, err = db.Prepare(`
			INSERT INTO translate(
				verb_id,
				ru_Ru
			) values(?,?)
		`)
		if err != nil {
			log.Fatal("Error prepare query in translate db. ", err)
		}
		_, err = stmt.Exec(id, list[6])
		if err != nil {
			log.Fatal("Error insert in db. ", err)
		}
	})
}
