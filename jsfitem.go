package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

const itemMode = "item"

type jsfItemArgs struct {
	siteURL       *string
	articlePath   *string
	title         *string
	contentFile   *string
	datePublished *string
	dateModified  *string
	tags          *string
	attachments   *string
}

func makeJSFItem() (jsfItemArgs, *flag.FlagSet) {
	var j jsfItemArgs
	f1 := flag.NewFlagSet(itemMode, flag.ContinueOnError)

	j.siteURL = f1.String("siteurl", "http://ratan.blog", "Your website URL")
	j.title = f1.String("title", "Untitled", "The title of the article")
	j.articlePath = f1.String("folder", "untitled", "The path of the article relative to your website")
	j.contentFile = f1.String("content", "content.html", "File with the article content")
	j.datePublished = f1.String("pdate", todayYYYYMMDD(), "Publish date (gregorian) in YYYY-MM-DD format")
	j.dateModified = f1.String("mdate", todayYYYYMMDD(), "Modify date (gregorian) in YYYY-MM-DD format")
	j.tags = f1.String("tags", "", "Comma-seperated tags")
	j.attachments = f1.String("attach", "", "Comma-seperated files to attach")

	return j, f1
}

type jsfAttachment struct {
	URL      string `json:"url"`
	MIMEType string `json:"mime_type"`
}

type jsfItem struct {
	ID            string          `json:"id"`
	URL           string          `json:"url"`
	Title         string          `json:"title"`
	ContentHTML   string          `json:"content_html"`
	DatePublished string          `json:"date_published"`
	DateModified  string          `json:"date_modified"`
	Tags          []string        `json:"tags"`
	Attachments   []jsfAttachment `json:"attachments"`
}

func buildItem(ja jsfItemArgs) {
	var ji jsfItem
	ji.Title = *(ja.title)
	ji.ID = path.Join(*(ja.siteURL), *(ja.articlePath))
	ji.URL = ji.ID
	ji.Tags = strings.Split(*(ja.tags), ",")

	var err error
	ji.DatePublished, err = jsfDate(*(ja.datePublished))
	killOnError(err)

	ji.DateModified, err = jsfDate(*(ja.dateModified))
	killOnError(err)

	aList := strings.Split(*(ja.attachments), ",")
	for _, attachName := range aList {
		if len(attachName) > 0 {
			a := buildAttachment(attachName, ji.URL)
			ji.Attachments = append(ji.Attachments, a)
		}
	}

	b0, err := ioutil.ReadFile(*(ja.contentFile))
	killOnError(err)

	ji.ContentHTML = string(b0)

	e := json.NewEncoder(os.Stdout)
	e.SetEscapeHTML(false)
	err = e.Encode(ji)
	killOnError(err)
}

func buildAttachment(filename string, baseURL string) jsfAttachment {
	var res jsfAttachment

	res.URL = path.Join(baseURL, filename)
	f, err := os.Open(filename)
	killOnError(err)

	b1 := make([]byte, 512)
	_, err = f.Read(b1)
	printOnError(err)

	err = f.Close()
	killOnError(err)

	res.MIMEType = http.DetectContentType(b1)
	return res
}
