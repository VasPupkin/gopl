// Package boltdb - provide interface to bolt database for xkcd utility
package boltdb

import (
	"fmt"
	"log"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	dbFile = "xkcd.db"
)

type Db struct {
	db *bolt.DB
}

// OpenDataBase - opens boltDB (database name is hardcoded)
func OpenDataBase() *Db {
	d := new(Db)
	//opens xkcb.db data file in current directory
	// it will be created if does not exist
	// option Timeout for prevent deadlock if data file already opened
	var err error
	d.db, err = bolt.Open(dbFile, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("boltdb OpenDataBase: %v\n", err)
	}
	return d
}

// CloseDataBase - closes opened DB
func (d *Db) CloseDataBase() {
	err := d.db.Close()
	if err != nil {
		log.Fatalf("boltdb CloseDataBase: %v\n", err)
	}
}

// SaveComic - saves information about comic into DB
func (d *Db) SaveComic(Num int, Date time.Time, Title, Transcription, ImageURL, AltName []byte) (err error) {
	n := []byte(strconv.Itoa(Num)) // single calculation
	// begin transactions
	tx, err := d.db.Begin(true)
	if err != nil {
		return fmt.Errorf("boltdb SaveComic: %v\n", err)
	}
	// if return some error do rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// create bucket for comic
	b, err := tx.CreateBucketIfNotExists(n)
	if err != nil {
		return fmt.Errorf("boltdb SaveComic: %v\n", err)
	}
	// marshal time date stamp // TODO: handle error here
	dt, _ := Date.MarshalBinary()
	// insert data into bucket
	err = b.Put([]byte("Date"), dt)
	if err != nil {
		return fmt.Errorf("Put Date: %s", err)
	}
	err = b.Put([]byte("Title"), Title)
	if err != nil {
		return fmt.Errorf("Put Title: %s", err)
	}
	err = b.Put([]byte("Transcription"), Transcription)
	if err != nil {
		return fmt.Errorf("Put Transcription: %s", err)
	}
	err = b.Put([]byte("ImageURL"), ImageURL)
	if err != nil {
		return fmt.Errorf("Put ImageURL: %s", err)
	}
	err = b.Put([]byte("AltName"), AltName)
	if err != nil {
		return fmt.Errorf("Put Title: %s", err)
	}
	// create or update bucket for last comic counter
	last, err := tx.CreateBucketIfNotExists([]byte("last"))
	if err != nil {
		return fmt.Errorf("create bucket: %s", err)
	}
	err = last.Put([]byte("Num"), n)
	if err != nil {
		return fmt.Errorf("Put Last: %s", err)
	}
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetComicInfo - retrieves comoc information from DB
func (d *Db) GetComicInfo(num int) (date time.Time, title, transcription, imageURL, altName []byte) {
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(strconv.Itoa(num)))
		if b == nil {
			return fmt.Errorf("bucket not exist")
		}
		err := date.UnmarshalBinary(b.Get([]byte("Date")))
		if err != nil {
			return fmt.Errorf("Get Date: %s", err)
		}
		title = b.Get([]byte("Title"))
		transcription = b.Get([]byte("Transcription"))
		imageURL = b.Get([]byte("ImageURL"))
		altName = b.Get([]byte("AltName"))
		return nil
	})
	if err != nil {
		log.Fatalf("boltdb GetComicInfo: %v\n", err)
	}
	return
}

// GetLastNum - return last number of comic available in DB
func (d *Db) GetLastNum() (last int) {
	tx, err := d.db.Begin(false)
	if err != nil {
		log.Fatalf("boltdb GetLastNum: %v\n", err)
	}
	b := tx.Bucket([]byte("last"))
	if b != nil {
		var err error
		last, err = strconv.Atoi(string(b.Get([]byte("Num"))))
		if err != nil {
			log.Fatalf("boltdb GetLastNum: %v\n", err)
		}
	}
	tx.Rollback()
	return
}

// CheckComicExist - check comic by number
func (d *Db) CheckComicExists(num int) (exist bool) {
	d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(strconv.Itoa(num)))
		exist = (b != nil)
		return nil
	})
	return
}
