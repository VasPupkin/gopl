package main

import (
	"fmt"
	database "gopl/ch4/xkcd/boltdb" // bolt DB
	"log"

	"gopl/ch4/xkcd/dwnldr"
	"gopl/ch4/xkcd/informer"
	// database "gopl/ch4/xkcd/sqlitedb" // sqlite3 DB
	"runtime"

	"github.com/gosuri/uiprogress"
)

var (
	version   = "" // from makefile
	buildtime = "" // from makefile
	altname   = "" // from makefile
)

const update = 200 // maximum requests by default

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // use all available cpu cores
	showInfo()
	db := database.OpenDataBase()
	fmt.Print("Checking web for updates...\n")
	inDb := db.GetLastNum()
	inWeb := dwnldr.FindLastComic()
	fmt.Print("Done\n")
	var limit int
	if inDb < inWeb {
		fmt.Printf("Internal Db has %d comic, need to download %d from web\n", inDb, inWeb-inDb)
		if inDb+update < inWeb {
			limit = update
		} else {
			limit = inWeb - inDb
		}
		fmt.Printf("Downloading next %d...\n", limit)
		// prepare progress bar
		bar := uiprogress.AddBar(limit).AppendCompleted().PrependElapsed()
		bar.PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("Task: (%0.3d/%0.3d)", b.Current(), limit)
		})
		uiprogress.Start()
		// TODO: Do it go routines
		for i := 0; i < limit; {
			i++
			if (i + inDb) == 404 {
				continue
			}
			c := dwnldr.GetComicInfo(i + inDb)
			err := db.SaveComic(c.Num, c.Date, []byte(c.Title), []byte(c.Transcription), []byte(c.ImageURL), []byte(c.AltName))
			if err != nil {
				log.Fatal(err)
			}
			bar.Incr()
		}
		uiprogress.Stop()
		fmt.Print("Download complete\n")
	} else {
		fmt.Printf("Internal Db is up to date, and contain %d comic\n", inDb)
	}
	db.CloseDataBase()
	informer.MainMenu(inDb + limit)
}

func showInfo() {
	if version != "" {
		fmt.Printf("Version: %s - %s\n", version, altname)
	}
	if buildtime != "" {
		fmt.Printf("Build time: %s\n", buildtime)
	}
}
