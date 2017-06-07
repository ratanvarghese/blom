package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/feeds"
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
const jfPath = "feeds/json"
const atomPath = "feeds/atom"
const rssPath = "feeds/rss"
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
	template := defaultTemplate
	blank := ""
	args.attach = &blank
	args.title = &blank
	args.tags = &blank
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
	jf.FeedUrl = path.Join(defaultHomePage, jfPath)
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
		curPath := jfPath
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
	templateFile := "../template.html"
	args.template = &templateFile

	newContent := fmt.Sprintf("%s<br /><a href=\"%s\">[Permalink]</a>", string(latestItem.ContentHTML), latestItem.URL)
	runTemplate(latestItem, args, newContent)
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
	articleE.Date = template.HTML("")
	articleE.Today = template.HTML(fmt.Sprintf("Today is %s.", dualDateFormat(time.Now().Format(time.RFC3339))))
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

func tagSort(itemList []jsfItem) (map[string][]jsfItem, []string) {
	res := make(map[string][]jsfItem)
	tagList := make([]string, 0)
	for _, ji := range itemList {
		for _, tag := range ji.Tags {
			if len(res[tag]) == 0 && len(tag) > 0 {
				tagList = append(tagList, tag)
			}
			res[tag] = append(res[tag], ji)
		}
	}
	sort.Strings(tagList)
	return res, tagList
}

func tagsPageLines(itemList []jsfItem) []string {
	outputLines := make([]string, 0)
	tagMap, tagList := tagSort(itemList)
	for _, tag := range tagList {
		outputLines = append(outputLines, fmt.Sprintf("<h3>%v</h3>", strings.Title(tag)))
		outputLines = append(outputLines, "<ul>")
		for _, ji := range tagMap[tag] {
			outputLines = append(outputLines, fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", ji.URL, ji.Title))
		}
		outputLines = append(outputLines, "</ul>")
	}
	return outputLines
}

func printTagsPage(itemList []jsfItem) {
	lineList := tagsPageLines(itemList)
	if err := os.Chdir("tags"); err != nil {
		panic(err)
	}

	var articleE articleExport
	articleE.Title = "Tags"
	articleE.Date = template.HTML("")
	articleE.Today = template.HTML(fmt.Sprintf("Today is %s.", dualDateFormat(time.Now().Format(time.RFC3339))))
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

func jfItemToGorilla(ji jsfItem) feeds.Item {
	var gi feeds.Item
	gi.Title = ji.Title
	gi.Link = &feeds.Link{Href: ji.URL}
	gi.Created, _ = time.Parse(time.RFC3339, ji.DatePublished)
	gi.Updated, _ = time.Parse(time.RFC3339, ji.DateModified)
	gi.Id = ji.URL
	return gi
}

func legacyFeeds(itemList []jsfItem) {
	var gf feeds.Feed
	gf.Title = defaultHomePage
	gf.Link = &feeds.Link{Href: siteURL}
	gf.Created = time.Now()

	gfItemList := make([]feeds.Item, len(itemList))
	gfItemPtrList := make([]*feeds.Item, len(itemList))
	for i, ji := range itemList {
		gfItemList[i] = jfItemToGorilla(ji)
		gfItemPtrList[i] = &(gfItemList[i])
	}
	gf.Items = gfItemPtrList
	atom, err := gf.ToAtom()
	if err != nil {
		panic(err)
	}

	rss, err := gf.ToRss()
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(atomPath, []byte(atom), 0664)
	ioutil.WriteFile(rssPath, []byte(rss), 0664)
}

func doUpdate() {
	jfItems := makeItemList()
	sort.Sort(byDatePublished(jfItems))

	if len(jfItems) > 0 {
		makeHomepage(jfItems[0])
	}
	paginatedPrint(jfItems)
	printArchive(jfItems)
	printTagsPage(jfItems)
	legacyFeeds(jfItems)
}
