package model

import (
	"bytes"
	_ "database/sql"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
	_ "github.com/mattn/go-sqlite3"
)

type Verb struct {
	Present         string
	PrTranscription string
	PastSimple      string
	PsTranscription string
	PastPartisiple  string
	PpTransription  string
	Translate       Translate
}

type Translate struct {
	language string
	Value    string
}

func (v *Verb) maxLength() int {
	vStrings := []string{v.Present, v.PrTranscription, v.PastSimple, v.PsTranscription, v.PastPartisiple, v.PpTransription}
	var max int
	for _, vString := range vStrings {
		if vLen := utf8.RuneCountInString(vString); max < vLen {
			max = vLen
		}
	}
	return max
}

func (v *Verb) getFormat() string {
	// length := strconv.FormatInt(int64(v.maxLength()), 10)
	buf := &bytes.Buffer{}
	for i := 1; i <= 3; i++ {
		buf.WriteString("%s")
		buf.WriteString(" ")
		buf.WriteString("%s")
		buf.WriteString(" -" + strconv.FormatInt(int64(i), 10) + "th. form")
		buf.WriteRune('\n')
	}
	return buf.String()
}

func (v *Verb) Find(c *AppDb, verb string, local string) error {
	searchString := strings.ToLower(verb)
	stmnt, err := c.db.Prepare(`
		SELECT 
		v.present,
		v.p_transcription,
		v.past_simple,
		v.ps_transcription,
		v.past_partisiple,
		v.pp_transription,
		tr.ru_Ru as translate FROM verbs as v
		LEFT JOIN 'translate' as tr ON v.id = tr.verb_id
		WHERE v.present = ?
		OR  v.past_simple= ?
		OR  v.past_partisiple= ?
		OR  translate LIKE ?
		LIMIT 1
	`)

	if err != nil {
		return fmt.Errorf("Error in find query %s", err.Error())
	}

	row := stmnt.QueryRow(searchString, searchString, searchString, "%"+searchString+"%")

	err = row.Scan(
		&v.Present,
		&v.PrTranscription,
		&v.PastSimple,
		&v.PsTranscription,
		&v.PastPartisiple,
		&v.PpTransription,
		&v.Translate.Value)
	if err != nil {
		return fmt.Errorf("I can find '" + verb + "' in my data set")
	}

	return nil
}

func (v *Verb) String() string {
	format := v.getFormat()
	return fmt.Sprintf(format,
		v.Present,
		v.PrTranscription,
		v.PastSimple,
		v.PsTranscription,
		v.PastPartisiple,
		v.PpTransription,
	)
}
