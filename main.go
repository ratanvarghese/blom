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

type templateArgs struct {
	applyTemplate *bool
	templateFile  *string
	contentFile   *string
	articleTitle  *string
	outputFile    *string
	styleSheet    *string
	date          *string
}

const gDateFormat = "2006-01-02"

func main() {
	var ta templateArgs
	ta.applyTemplate = flag.Bool("apply", false, "Apply a template")
	ta.templateFile = flag.String("template", "template.html", "File to use as a template")
	ta.contentFile = flag.String("content", "content.html", "File with HTML content to insert in template")
	ta.articleTitle = flag.String("title", "Untitled", "Title of the article")
	ta.outputFile = flag.String("output", "index.html", "Filename of output of the executed template")
	ta.styleSheet = flag.String("style", "../style.css", "Filename of stylesheet")
	ta.date = flag.String("date", time.Now().Format(gDateFormat), "Gregorian date in format YYYY-MM-DD, defaults to today")

	flag.Parse()

	runTemplate(ta)
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
	if !*(ta.applyTemplate) {
		return
	}
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
