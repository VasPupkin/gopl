package informer

import "text/template"

const (
	templ = `
Comic #{{.Num}}
----------------------------------------------------------------------
Title:  {{.Title}}
Date:   {{.Date | dateString}}
Transcription: {{.Transcription | formatString}}
Altern Image:  {{.AltName | formatString}}
Image link: {{.ImageURL}}
----------------------------------------------------------------------

`

	tLen    = 70 // template length
	tIndent = 15 // first linr indent (Transcription: )
)

/*
<--                             70                                 -->
*/

var comicInfo = template.Must(template.New("comicInfo").
	Funcs(template.FuncMap{"dateString": dateString,
		"formatString": formatString}).Parse(templ))
