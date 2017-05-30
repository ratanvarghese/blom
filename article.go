package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const articleMode = "article"
const templateFile = "template.html"
const siteURL = "http://ratan.blog"
const contentFile = "content.html"
const itemFile = "item.json"
const listSeperator = ","
const outputWebpage = "index.html"

type articleArgs struct {
	attach     *string
	title      *string
	tags       *string
	localstyle *bool
}

type jsfAttachment struct {
	URL      string `json:"url"`
	MIMEType string `json:"mime_type"`
	valid    bool
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

type articleExport struct {
	Title       string
	Stylesheet  string
	Date        string
	Today       string
	ContentHTML template.HTML
}

func makeArticleArgs() (articleArgs, *flag.FlagSet) {
	var args articleArgs
	fset := flag.NewFlagSet(articleMode, flag.ContinueOnError)

	args.attach = fset.String("attach", "", "Comma-seperated files to attach")
	args.title = fset.String("title", "", "Title of the article")
	args.tags = fset.String("tags", "", "Comma-seperated tags")
	args.localstyle = fset.Bool("localstyle", false, "Use stylesheet in same folder as output")

	return args, fset
}

func argsFromItemFile() (string, string, string, string) {
	datePublished := time.Now().Format(time.RFC3339)

	prevFileContent, err := ioutil.ReadFile(itemFile)
	if err != nil {
		msg := strings.ToLower(err.Error())
		if !strings.Contains(msg, "no such file") {
			log.Print(err)
		}
		return "", "", "", datePublished
	}

	var ji jsfItem
	err = json.Unmarshal(prevFileContent, &ji)
	if err != nil {
		log.Print(err)
	}

	attach := ""
	for _, a := range ji.Attachments {
		attach = strings.Join([]string{attach, filepath.Base(a.URL)}, listSeperator)
	}
	title := ji.Title
	tags := (strings.Join(ji.Tags, listSeperator))
	datePublished = ji.DatePublished
	return attach, title, tags, datePublished
}

func curDir() string {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return filepath.Base(p)
}

func makeItem(args articleArgs, datePublished string, content string) jsfItem {
	var res jsfItem
	res.Title = *(args.title)
	res.Tags = strings.Split(*(args.tags), listSeperator)
	res.URL = path.Join(siteURL, curDir())
	res.ID = res.URL
	res.DatePublished = datePublished
	res.ContentHTML = content

	attachList := strings.Split(*(args.attach), listSeperator)
	for _, attachName := range attachList {
		if len(attachName) > 0 {
			a := buildAttachment(attachName, siteURL)
			if a.valid {
				res.Attachments = append(res.Attachments, a)
			}
		}
	}

	info, err := os.Stat(contentFile)
	if err != nil {
		panic(err)
	}
	res.DateModified = info.ModTime().Format(time.RFC3339)
	if strings.Compare(res.DateModified, res.DatePublished) < 0 {
		res.DateModified = res.DatePublished
		//As far as the outside world is concerned, the article does not exist
		//before it is published.
	}

	f, err := os.Create(itemFile)
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.Encode(res)
	if err != nil {
		panic(err)
	}
	return res
}

func buildAttachment(filename string, baseURL string) jsfAttachment {
	res := jsfAttachment{"", "", false}

	f, err := os.Open(filename)
	if err != nil {
		log.Print(err)
		return res
	}

	b1 := make([]byte, 512)
	_, err = f.Read(b1)
	if err != nil {
		log.Print(err)
		return res
	}
	err = f.Close()
	if err != nil {
		log.Print(err)
	}

	res.URL = path.Join(baseURL, filename)
	res.MIMEType = http.DetectContentType(b1)
	res.valid = true
	return res
}

func dualDateFormat(RFCDate string) string {
	gDate, err := time.Parse(time.RFC3339, RFCDate)
	if err != nil {
		panic(err)
	}

	const outputGDateFormat = "Monday, 2 January, 2006 CE"
	tqDate := tqtime.LongDate(gDate.Year(), gDate.YearDay())
	tqDateBetter := strings.Replace(tqDate, "After Tranquility", "AT", 1)
	gDateStr := gDate.Format(outputGDateFormat)
	return fmt.Sprintf("%s [Gregorian: %s]", tqDateBetter, gDateStr)
}

func runTemplate2(ji jsfItem, args articleArgs, content string) {
	var articleE articleExport
	articleE.Title = ji.Title
	if *(args.localstyle) {
		articleE.Stylesheet = "style.css"
	} else {
		articleE.Stylesheet = "../style.css"
	}
	articleE.Date = dualDateFormat(ji.DateModified)
	articleE.Today = fmt.Sprintf("Today is %s.", dualDateFormat(time.Now().Format(time.RFC3339)))
	articleE.ContentHTML = template.HTML(content)

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(outputWebpage)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, articleE)
	if err != nil {
		panic(err)
	}
}

func buildArticle(args articleArgs) {
	articleContent, err := ioutil.ReadFile(contentFile)
	if err != nil {
		panic(err)
	}

	oldAttach, oldTitle, oldTags, datePublished := argsFromItemFile()
	if len(*args.attach) < 1 {
		args.attach = &oldAttach
	}
	if len(*args.title) < 1 {
		args.title = &oldTitle
	}
	if len(*args.tags) < 1 {
		args.tags = &oldTags
	}

	ji := makeItem(args, datePublished, string(articleContent))
	runTemplate2(ji, args, string(articleContent))
}
