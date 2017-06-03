package main

import (
	"encoding/json"
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

const updateMode = "update"
const defaultVersion = "https://jsonfeed.org/version/1"
const defaultHomePage = "ratan.blog"
const feedPath = "feeds/json"
const pageLen = 15

type jsfMain struct {
	Version     string    `json:"version"`
	Title       string    `json:"title"`
	HomePageUrl string    `json:"home_page_url"`
	FeedUrl     string    `json:"feed_url"`
	NextUrl     string    `json:"next_url,omitempty"`
	Items       []jsfItem `json:"items"`
}

type byDatePublished []jsfItem

func (b byDatePublished) Len() int {
	return len(b)
}

func (b byDatePublished) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byDatePublished) Less(i, j int) bool {
	ti, _ := time.Parse(time.RFC3339, b[i].DatePublished)
	tj, _ := time.Parse(time.RFC3339, b[j].DatePublished)
	return ti.After(tj)
}

func noFrillsArticle() articleArgs {
	var args articleArgs
	stylesheet := defaultStylesheet
	template := defaultTemplate
	blank := ""
	args.attach = &blank
	args.title = &blank
	args.tags = &blank
	args.style = &stylesheet
	args.template = &template
	return args
}

func makeItemList() []jsfItem {
	jiList := make([]jsfItem, 0)
	dir, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, folder := range dir {
		itemInFolder := path.Join(folder.Name(), itemFile)
		if _, err := os.Stat(itemInFolder); os.IsNotExist(err) {
			continue
		}
		if err := os.Chdir(folder.Name()); err != nil {
			continue
		}
		args := noFrillsArticle()
		buildArticle(args)
		if err := os.Chdir(".."); err != nil {
			panic(err)
		}
		itemFileContent, err := ioutil.ReadFile(itemInFolder)
		if err != nil {
			panic(err)
		}
		var ji jsfItem
		err = json.Unmarshal(itemFileContent, &ji)
		if err != nil {
			panic(err)
		}
		jiList = append(jiList, ji)

	}
	return jiList
}

func defaultJsfMain() jsfMain {
	var jf jsfMain
	jf.Version = defaultVersion
	jf.HomePageUrl = defaultHomePage
	jf.Title = defaultHomePage
	jf.FeedUrl = path.Join(defaultHomePage, feedPath)
	return jf
}

func paginatedPrint(itemList []jsfItem) {
	jf := defaultJsfMain()
	listLen := len(itemList)
	for i := 0; i < listLen; i += pageLen {
		pageNum := i / pageLen
		pageEnd := i + pageLen
		if listLen >= pageEnd {
			jf.NextUrl = fmt.Sprintf("%v%v", jf.FeedUrl, pageNum+1)
		} else {
			pageEnd = listLen
		}
		jf.Items = itemList[i:pageEnd]
		curPath := feedPath
		if pageNum > 0 {
			curPath = fmt.Sprintf("%v%v", curPath, pageNum)
		}
		f, err := os.Create(curPath)
		if err != nil {
			panic(err)
		}
		enc := json.NewEncoder(f)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "\t")
		enc.Encode(jf)
	}
}

func makeHomepage(latestItem jsfItem) {
	args := noFrillsArticle()
	style := "style.css"
	template := "../template.html"
	args.style = &style
	args.template = &template
	runTemplate(latestItem, args, latestItem.ContentHTML)
}

func archiveLines(itemList []jsfItem) []string {
	var gt1 time.Time
	outputLines := make([]string, 0)
	for i, ji := range itemList {
		g1Year := gt1.Year()
		g1YearDay := gt1.YearDay()
		tq1Year := tqtime.Year(g1Year, g1YearDay)
		tq1Mon := tqtime.Month(g1Year, g1YearDay)
		tq1Day := tqtime.Day(g1Year, g1YearDay)

		gt2, _ := time.Parse(time.RFC3339, ji.DatePublished)
		g2Year := gt2.Year()
		g2YearDay := gt2.YearDay()
		tq2Year := tqtime.Year(g2Year, g2YearDay)
		tq2Mon := tqtime.Month(g2Year, g2YearDay)
		tq2Day := tqtime.Day(g2Year, g2YearDay)

		isSpecialDay := (tq2Mon == tqtime.SpecialDay)

		if (tq1Year != tq2Year) || (tq1Mon != tq2Mon) || (isSpecialDay && (tq1Day != tq2Day)) {
			if !gt1.IsZero() {
				outputLines = append(outputLines, "<ul>")
			}
			if isSpecialDay {
				outputLines = append(outputLines, fmt.Sprintf("<h3>%v, %v AT</h3>", tqtime.DayName(tq2Day), tq2Year))
			} else {
				outputLines = append(outputLines, fmt.Sprintf("<h3>%v, %v AT</h3>", tq2Mon.String(), tq2Year))
			}
			outputLines = append(outputLines, "<ul>")
		}
		outputLines = append(outputLines, fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", ji.URL, ji.Title))
		gt1 = gt2
		if i == (len(itemList) - 1) {
			outputLines = append(outputLines, "</ul>")
		}
	}
	return outputLines
}

func printArchive(itemList []jsfItem) {
	lineList := archiveLines(itemList)
	if err := os.Chdir("archive"); err != nil {
		panic(err)
	}

	var articleE articleExport
	articleE.Title = "Archive"
	articleE.Stylesheet = "../style.css"
	articleE.Date = ""
	articleE.Today = fmt.Sprintf("Today is %s.", dualDateFormat(time.Now().Format(time.RFC3339)))
	articleE.ContentHTML = template.HTML(strings.Join(lineList, "\n"))

	tmpl, err := template.ParseFiles("../../template.html")
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

	if err := os.Chdir(".."); err != nil {
		panic(err)
	}
}

func doUpdate() {
	jfItems := makeItemList()
	sort.Sort(byDatePublished(jfItems))

	if len(jfItems) > 0 {
		makeHomepage(jfItems[0])
	}
	paginatedPrint(jfItems)
	printArchive(jfItems)
}
