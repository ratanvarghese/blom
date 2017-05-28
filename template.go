package main

import (
	"bufio"
	"flag"
	"html/template"
	"io/ioutil"
	"os"
)

const templateMode = "template"

type templateArgs struct {
	templateFile *string
	contentFile  *string
	articleTitle *string
	styleSheet   *string
	date         *string
}

func makeTemplateArgs() (templateArgs, *flag.FlagSet) {
	var ta templateArgs
	f1 := flag.NewFlagSet(templateMode, flag.ContinueOnError)

	ta.templateFile = f1.String("template", "template.html", "File to use as a template")
	ta.contentFile = f1.String("content", "content.html", "File with HTML content to insert in template")
	ta.articleTitle = f1.String("title", "Untitled", "Title of the article")
	ta.styleSheet = f1.String("style", "../style.css", "Filename of stylesheet")
	ta.date = f1.String("date", todayYYYYMMDD(), "Gregorian date in format YYYY-MM-DD, defaults to today")

	return ta, f1
}

type articleExport struct {
	Title       string
	Stylesheet  string
	Date        string
	Today       string
	ContentHTML template.HTML
}

func runTemplate(ta templateArgs) {
	t, err := template.ParseFiles(*(ta.templateFile))
	killOnError(err)

	content, err := ioutil.ReadFile(*(ta.contentFile))
	killOnError(err)

	var articleE articleExport
	articleE.Title = *(ta.articleTitle)
	articleE.Stylesheet = *(ta.styleSheet)
	articleE.Date, err = webpageDate(*(ta.date))
	articleE.Today = headerDate()
	killOnError(err)

	w := bufio.NewWriter(os.Stdout)
	articleE.ContentHTML = template.HTML(content)

	err = t.Execute(w, articleE)
	printOnError(err)

	err = w.Flush()
	printOnError(err)
}
