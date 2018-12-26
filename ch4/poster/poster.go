package main

import (
	"flag"
	"fmt"
	"gopl/ch4/poster/omdbapi"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var (
	title     = flag.String("title", "", "Movie title")
	year      = flag.String("year", "", "Movie release year")
	version   = "" // from makefile
	buildtime = "" // from makefile
	altname   = "" // from makefile
)

func main() {
	showInfo()
	flag.Parse()
	if *title == "" {
		log.Fatalf("Title can not be a empty\n")
	}
	if *year != "" {
		var y int
		var err error
		if y, err = strconv.Atoi(*year); err != nil {
			log.Fatalf("Wrong year value\n")
		}
		if y < 1900 || y > time.Now().Year() {
			log.Fatalf("Wrong year value\n")
		}
	}
	movie, ok := omdbapi.FindMovie(*title, *year)
	// if not found or other erors
	if !ok {
		fmt.Printf("Not Found movie: %s\n", *title)
		os.Exit(0)
	}
	fmt.Printf("Found movie!\nTitle: %s\nYear: %s\n", movie.Title, movie.Year)
	if movie.Poster == "N/A" || movie.Poster == "" {
		fmt.Printf("Poster for %s unavailable\n", movie.Title)
		os.Exit(0)
	}
	fmt.Printf("Show %s poster? (y/n) ", movie.Title)
	// ask to show poster
	var inp string
	fmt.Scanln(&inp)
	if inp == "n" || inp == "N" {
		fmt.Printf("Canceled\n")
		os.Exit(0)
	}
	// if yes  show it
	showPoster(movie)
}

// prints invormatin
func showInfo() {
	if version != "" {
		fmt.Printf("Version: %s - %s\n", version, altname)
	}
	if buildtime != "" {
		fmt.Printf("Build time: %s\n", buildtime)
	}
}

func showPoster(m *omdbapi.Movie) {
	poster, ok := m.DownloadPoster()
	//  if can`t download or other errors
	if !ok {
		fmt.Printf("Download Error\n")
		os.Exit(0)
	}
	// make temp.file
	fname := fmt.Sprintf("img%s", poster.Type)
	if err := ioutil.WriteFile(fname, poster.Image, 0644); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(fname) // clean up
	var cmd *exec.Cmd
	// select image viewer
	if isCommandAvailable("feh") {
		cmd = exec.Command("feh", "-ZxF", fname)
	} else {
		cmd = exec.Command("xdg_open", fname)
	}
	//wait until app will be closed
	_, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
