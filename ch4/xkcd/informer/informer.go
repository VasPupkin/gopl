package informer

import (
	"bytes"
	"fmt"
	database "gopl/ch4/xkcd/boltdb" // boltdb
	// database "gopl/ch4/xkcd/sqlitedb" //sqlite3 db
	"log"
	"os"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"
)

var (
	spc = []byte{32}
	lf  = []byte{10}
)

func MainMenu(num int) {
	for {
		fmt.Printf("Input comic number [1 - %d] or [Q/q] for exit:", num)
		inp := ""
		fmt.Scanln(&inp)
		if inp == "q" || inp == "Q" || inp == "" {
			fmt.Println("Exiting")
			break
		}
		comic, err := strconv.Atoi(inp)
		if err != nil {
			fmt.Printf("Wrong input: %d\n", inp)
			continue
		}
		if comic > num || comic < 1 {
			fmt.Printf("Wrong comic number: %d\n", comic)
			continue
		}
		showComicInfo(comic)
	}
}

func showComicInfo(num int) {
	type comic struct {
		Num           int
		Date          time.Time
		Title         string
		Transcription []byte
		ImageURL      string
		AltName       []byte
	}
	c := comic{}
	d := database.OpenDataBase()
	defer d.CloseDataBase()
	// check that the comic truly exists in DB
	if !d.CheckComicExists(num) {
		fmt.Printf("Unavailable comic #%d\n", num)
		return
	}
	date, title, transcription, imageURL, altName := d.GetComicInfo(num)
	c.Date = date
	c.Num = num
	c.Title = string(title)
	c.Transcription = transcription
	c.ImageURL = string(imageURL)
	c.AltName = altName
	if err := comicInfo.Execute(os.Stdout, c); err != nil {
		log.Fatalf("Template execution: %v\n", err)
	}
}

func dateString(t time.Time) string {
	month := []string{"Jan", "Feb", "Mar", "Apr", "May",
		"Jun", "Jul", "Aug", "Sen", "Oct", "Nov", "Dec"}
	y, m, d := t.Date()
	return fmt.Sprintf("%d %s %d", d, month[int(m)-1], y)
}

func formatString(in []byte) string {
	if len(in) == 0 {
		return "None"
	}
	// remove all line feeds and leading/ending spaces
	in = bytes.Replace(in, lf, spc, -1)
	in = bytes.TrimPrefix(in, spc)
	in = bytes.TrimSuffix(in, spc)
	// replace all duplicated spaces and split to words
	words := bytes.Split(spacer(spaceReplacer(in)), spc)
	strLen := tIndent // first indent
	out := ""
	for _, word := range words {
		wLen := utf8.RuneCount(word)
		if (strLen + wLen) > tLen {
			out += "\n"
			strLen = 0 //new line
		}
		out += string(word)
		out += " "
		strLen += (wLen + 1)
	}
	return out
}

// remove duplicated spases in string
func spacer(str []byte) []byte {
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRune(str[i:])
		if unicode.IsSpace(r) {
			if (i + size) < len(str) {
				if r, _ := utf8.DecodeRune(str[i+size:]); unicode.IsSpace(r) {
					copy(str[i:], str[i+size:])
					str = spacer(str[:len(str)-size])
				}
			}
		}
		i += size
	}

	return str
}

// replaces any kind of the space to ASCII space
func spaceReplacer(str []byte) []byte {
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRune(str[i:])
		if unicode.IsSpace(r) {
			copy(str[i:i+1], []byte{32})
		}
		i += size
	}
	return str
}
