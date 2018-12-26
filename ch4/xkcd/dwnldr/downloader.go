// download and update xkcd comic database
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

type Comic struct {
	Num           int
	Date          time.Time
	Title         string `json:"safe_title"`
	Transcription string `json:"transcript"`
	ImageURL      string `json:"img"`
	AltName       string `json:"alt"`
}

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
