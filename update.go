package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
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

func noFrillsBuildArticle() {
	var args articleArgs
	stylesheet := defaultStylesheet
	template := defaultTemplate
	blank := ""
	args.attach = &blank
	args.title = &blank
	args.tags = &blank
	args.style = &stylesheet
	args.template = &template
	buildArticle(args)
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
		noFrillsBuildArticle()
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

func doUpdate() {
	jfItems := makeItemList()
	sort.Sort(byDatePublished(jfItems))
	paginatedPrint(jfItems)
}
