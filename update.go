package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/feeds"
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
const defaultPageLen = 15

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
	jf.HomePageURL = defaultHomePage
	jf.Title = defaultHomePage
	jf.FeedURL = path.Join(defaultHomePage, jsfPath)
	return jf
}

func paginatedPrint(itemList []jsfItem) {
	jf := defaultJsfMain()
	listLen := len(itemList)
	for i := 0; i < listLen; i += defaultPageLen {
		pageNum := i / defaultPageLen
		pageEnd := i + defaultPageLen
		if listLen >= pageEnd {
			jf.NextURL = fmt.Sprintf("%v%v", jf.FeedURL, pageNum+1)
		} else {
			pageEnd = listLen
		}
		jf.Items = itemList[i:pageEnd]
		curPath := jsfPath
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
		err = enc.Encode(jf)
		if err != nil {
			panic(err)
		}
	}
}

func makeHomepage(latestItem jsfItem) {
	args := noFrillsArticle()
	templateFile := "../template.html"
	args.template = &templateFile

	newContent := fmt.Sprintf("%s<br /><a href=\"%s\">[Permalink]</a>", latestItem.ContentHTML, latestItem.URL)
	runTemplate(latestItem, args, newContent)
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
	gi.Description = ji.ContentHTML
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

	err = ioutil.WriteFile(atomPath, []byte(atom), 0664)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(rssPath, []byte(rss), 0664)
	if err != nil {
		panic(err)
	}
}

func doUpdate() {
	jfItems := makeItemList()
	sort.Sort(byPublishedDescend(jfItems))

	if len(jfItems) > 0 {
		makeHomepage(jfItems[0])
	}
	paginatedPrint(jfItems)
	printArchive(jfItems)
	printTagsPage(jfItems)
	legacyFeeds(jfItems)
}
