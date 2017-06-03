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
	"net/url"
	"os"
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

const defaultStylesheet = "../style.css"
const defaultTemplate = "../../template.html"

type articleArgs struct {
	attach   *string
	title    *string
	tags     *string
	style    *string
	template *string
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
	Date        template.HTML
	Today       template.HTML
	ContentHTML template.HTML
}

func makeArticleArgs() (articleArgs, *flag.FlagSet) {
	var args articleArgs
	fset := flag.NewFlagSet(articleMode, flag.ContinueOnError)

	args.attach = fset.String("attach", "", "Comma-seperated files to attach")
	args.title = fset.String("title", "", "Title of the article")
	args.tags = fset.String("tags", "", "Comma-seperated tags")
	args.style = fset.String("style", defaultStylesheet, "Filename of stylesheet")
	args.template = fset.String("template", defaultTemplate, "Filename of template file")

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

	base, err := url.Parse(siteURL)
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(curDir())
	if err != nil {
		panic(err)
	}

	res.URL = base.ResolveReference(u).String()
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
	tMod := info.ModTime()
	tPub, err := time.Parse(time.RFC3339, res.DatePublished)
	if err != nil {
		panic(err)
	}
	if tMod.Before(tPub) {
		res.DateModified = res.DatePublished
	} else {
		res.DateModified = info.ModTime().Format(time.RFC3339)
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

	base, err := url.Parse(baseURL)
	if err != nil {
		log.Print(err)
		return res
	}

	u, err := url.Parse(filename)
	if err != nil {
		log.Print(err)
		return res
	}

	res.URL = base.ResolveReference(u).String()
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
	return fmt.Sprintf("%s<br />[Gregorian: %s]", tqDateBetter, gDateStr)
}

func runTemplate(ji jsfItem, args articleArgs, content string) {
	var articleE articleExport
	articleE.Title = ji.Title
	articleE.Stylesheet = *(args.style)
	articleE.Date = template.HTML(dualDateFormat(ji.DatePublished))
	articleE.Today = template.HTML(fmt.Sprintf("Today is %s.", dualDateFormat(time.Now().Format(time.RFC3339))))
	articleE.ContentHTML = template.HTML(content)

	tmpl, err := template.ParseFiles(*(args.template))
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
	runTemplate(ji, args, string(articleContent))
}
