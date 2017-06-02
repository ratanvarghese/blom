package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"time"
)

const updateMode = "update"
const defaultVersion = "https://jsonfeed.org/version/1"
const defaultHomePage = "ratan.blog"
const defaultFeedUrl = "ratan.blog/feeds/json"

type jsfMain struct {
	Version     string    `json:"version"`
	Title       string    `json:"title"`
	HomePageUrl string    `json:"home_page_url"`
	FeedUrl     string    `json:"feed_url"`
	NextUrl     string    `json:"next_url"`
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

func doUpdate() {
	var jf jsfMain
	jf.Version = defaultVersion
	jf.HomePageUrl = defaultHomePage
	jf.Title = defaultHomePage
	jf.FeedUrl = defaultFeedUrl
	jf.Items = makeItemList()
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	sort.Sort(byDatePublished(jf.Items))
	enc.Encode(jf)
}
