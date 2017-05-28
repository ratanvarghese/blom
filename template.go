package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const templateMode = "template"
const gDateFormat = "2006-01-02"

type templateArgs struct {
	templateFile *string
	contentFile  *string
	articleTitle *string
	outputFile   *string
	styleSheet   *string
	date         *string
}

func makeTemplateArgs() (templateArgs, *flag.FlagSet) {
	var ta templateArgs
	f1 := flag.NewFlagSet("template", flag.ContinueOnError)

	ta.templateFile = f1.String("template", "template.html", "File to use as a template")
	ta.contentFile = f1.String("content", "content.html", "File with HTML content to insert in template")
	ta.articleTitle = f1.String("title", "Untitled", "Title of the article")
	ta.outputFile = f1.String("output", "index.html", "Filename of output of the executed template")
	ta.styleSheet = f1.String("style", "../style.css", "Filename of stylesheet")
	ta.date = f1.String("date", time.Now().Format(gDateFormat), "Gregorian date in format YYYY-MM-DD, defaults to today")

	return ta, f1
}

type articleExport struct {
	Title       string
	Stylesheet  string
	Date        string
	ContentHTML template.HTML
}

func killOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printOnError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func runTemplate(ta templateArgs) {
	t, err := template.ParseFiles(*(ta.templateFile))
	killOnError(err)

	content, err := ioutil.ReadFile(*(ta.contentFile))
	killOnError(err)

	f, err := os.Create(*(ta.outputFile))
	killOnError(err)

	var articleE articleExport
	articleE.Title = *(ta.articleTitle)
	articleE.Stylesheet = *(ta.styleSheet)

	const outputGDateFormat = "Monday, 2 January, 2006 CE"
	gDate, err := time.Parse(gDateFormat, *(ta.date))
	killOnError(err)

	tqDate := tqtime.LongDate(gDate.Year(), gDate.YearDay())
	tqDateBetter := strings.Replace(tqDate, "After Tranquility", "AT", 1)
	gDateStr := gDate.Format(outputGDateFormat)
	articleE.Date = fmt.Sprintf("%s [Gregorian: %s]", tqDateBetter, gDateStr)

	w := bufio.NewWriter(f)
	articleE.ContentHTML = template.HTML(content)

	err = t.Execute(w, articleE)
	printOnError(err)

	err = w.Flush()
	printOnError(err)

	err = f.Close()
	printOnError(err)
}
