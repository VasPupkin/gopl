//Package dwnldr - download and update xkcd comic database
package dwnldr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	url     = "https://xkcd.com/"
	srchStr = "Permanent link to this comic: https://xkcd.com/"
)

// Comic structure about comic
type Comic struct {
	//Num is comic number from xkcd.com,
	//in recieved json fild name "num"
	Num int
	//Conversion of sum three of fields: "day",
	// "month" and "year"
	Date time.Time
	// Recieved as "safe_title"
	Title string `json:"safe_title"`
	// Resieved as "transcrrept"
	Transcription string `json:"transcript"`
	// Resieved as "img", images are not downloaded
	ImageURL string `json:"img"`
	// Recieved as "alt", may be empty
	AltName string `json:"alt"`
}

// FindLastComic - request from web number of last published comic
func FindLastComic() (num int) {
	var permLink string
	resp, err := http.Get(url)
	if err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("request failed %s", resp.Status)
	}
	// close response anyway

	str := bufio.NewScanner(resp.Body)
	for str.Scan() {
		if strings.Contains(str.Text(), srchStr) {
			permLink = str.Text()
			break
		}
	}
	if len(permLink) > 0 {
		num, err = strconv.Atoi(strings.TrimRight(
			strings.TrimLeft(permLink, srchStr), `/<br />`))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Not found permanent link")
	}

	return
}

// GetComic info - downloads information aboit comic from web
func GetComicInfo(num int) *Comic {
	type date struct {
		Day   string
		Month string
		Year  string
	}
	cnv := func(s string) int {
		n, err := strconv.Atoi(s)
		if err != nil {
			log.Fatalf("Conversion: %v", err)
		}
		return n
	}
	url := fmt.Sprintf(`https://xkcd.com/%d/info.0.json`, num)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("request failed %s", resp.Status)
	}
	// str := bufio.NewScanner(resp.Body)
	// for str.Scan() {
	// 	fmt.Printf("%s\n", str.Text())
	// }
	var result Comic
	var d date
	r := struct {
		*Comic
		*date
	}{&result, &d}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&r); err != nil {
		log.Fatalf("Decode: %v", err)
	}
	result.Date, err = time.Parse("01-02-2006",
		fmt.Sprintf("%0.2d-%0.2d-%0.4d", cnv(d.Month), cnv(d.Day), cnv(d.Year)))
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}
	return &result
}
