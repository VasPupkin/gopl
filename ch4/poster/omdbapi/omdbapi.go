package omdbapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

var apiKey = "" // from make file

type Movie struct {
	Title  string
	Year   string
	Poster string
}

type PosterImage struct {
	Type  string
	Image []byte
}

func FindMovie(title string, year string) (*Movie, bool) {
	var req string
	// prepare request string
	if year == "" {
		req = fmt.Sprintf(`http://www.omdbapi.com/?apikey=%s&t=%s`,
			apiKey, url.QueryEscape(title))
	} else {
		req = fmt.Sprintf(`http://www.omdbapi.com/?apikey=%s&t=%s&y=%s`,
			apiKey, url.QueryEscape(title), year)

	}
	// get json from omdbapi.com
	resp, err := http.Get(req)
	if err != nil {
		log.Fatalf("omdbapi FindMovie get: %v\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("request failed %s", resp.Status)
	}
	var movie Movie
	rslt := struct {
		*Movie
		Response string
	}{&movie, ""}
	// partially decode json
	if err := json.NewDecoder(resp.Body).Decode(&rslt); err != nil {
		return &movie, false
	}
	return &movie, rslt.Response == "True"
}

func (m *Movie) DownloadPoster() (*PosterImage, bool) {
	pi := new(PosterImage)
	pi.Type = filepath.Ext(m.Poster)
	resp, err := http.Get(m.Poster)
	if err != nil {
		log.Fatalf("omdbapi DownloadPoster get: %v\n", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("request failed %s", resp.Status)
	}
	pi.Image, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Body read: %v\n", err)
		return pi, false
	}
	return pi, true
}
