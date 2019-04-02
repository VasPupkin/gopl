//Package sqlitedb - provide interface to sqlite database for xkcd utility
package sqlitedb

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "xkcd.sql3"
)

type Db struct {
	db *sql.DB
}

// SaveComic - saves information about comic into DB
func (d *Db) SaveComic(Num int, Date time.Time, Title, Transcription, ImageURL, AltName []byte) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("db save: %v\n", err)
	}
	stmt, err := tx.Prepare(
		"INSERT INTO main(Num, Date, Title, Transcription, ImageURL, AltName) values(?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("db save stmt: %v\n", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(Num, Date, Title, Transcription, ImageURL, AltName)
	if err != nil {
		return fmt.Errorf("db save exec: %v\n", err)
	}
	tx.Commit()
}

func (d *Db) GetComicInfo(num int) (date time.Time, title, transcription, imageURL, altName []byte) {
	stmt := `SELECT Date, Title, Transcription, ImageURL, AltName FROM main WHERE Num = ?`
	var tmpTime []byte
	err := d.db.QueryRow(stmt, num).Scan(&tmpTime, &title, &transcription, &imageURL, &altName)
	if err != nil {
		log.Fatalf("DB GetComicInfo: %v\n", err)
	}
	date, _ = time.Parse("2006-01-02", strings.TrimSuffix(string(tmpTime), " 00:00:00+00:00"))
	return
}

// OpenDataBase - opens sqlite database (database name is hardcoded)
func OpenDataBase() *Db {
	if _, err := os.Stat(dbFile); err != nil {
		if os.IsNotExist(err) {
			createNew()
		}
	}
	d := new(Db)
	var err error
	d.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("database OpenDataBase: %v\n", err)
	}
	return d
}

// CloseDataBase - closes opened DB
func (d *Db) CloseDataBase() {
	err := d.db.Close()
	if err != nil {
		log.Fatalf("DB close: %v\n", err)
	}
}

// GetLastNum - return last number of comic available in DB
func (d *Db) GetLastNum() int {
	cstmt := `SELECT MAX(Num) AS max FROM main`
	var num int
	err := d.db.QueryRow(cstmt).Scan(&num)
	if err != nil {
		log.Fatalf("DB GetlastNum: %v\n", err)
	}
	return num
}

// CheckComicExist - check comic by number
func (d *Db) CheckComicExists(num int) bool {
	var count int
	cstmt := `SELECT COUNT (Num) AS count FROM main WHERE Num = ?`
	err := d.db.QueryRow(cstmt, num).Scan(&count)
	if err != nil {
		log.Fatalf("DB CheckComicExists: %v\n", err)
	}
	return count == 1
}

func createNew() {
	db, err := sql.Open("sqlite3", dbFile)
	stmt := `
CREATE TABLE "main" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"Num"	INTEGER,
	"Date"	BLOB,
	"Title"	BLOB,
	"Transcription"	BLOB,
	"ImageURL"	BLOB,
	"AltName"	BLOB
)`
	stmt2 := `INSERT INTO "main" VALUES (0,0,0x00,0x00,0x00,0x00,0x00)`
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatalf("DB create New: %v\n", err)
	}
	_, err = db.Exec(stmt2)
	if err != nil {
		log.Fatalf("DB inser start values: %v\n", err)
	}
	db.Close()
}
