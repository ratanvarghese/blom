package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const jsfVersion = "https://jsonfeed.org/version/1"
const jsfPath = "feeds/json"
const atomPath = "feeds/atom"
const rssPath = "feeds/rss"

type jsfMain struct {
	Version     string    `json:"version"`
	Title       string    `json:"title"`
	HomePageURL string    `json:"home_page_url"`
	FeedURL     string    `json:"feed_url"`
	NextURL     string    `json:"next_url,omitempty"`
	Items       []jsfItem `json:"items"`
}

type jsfItemErr struct {
	item jsfItem
	err  error
}

type byPublishedDescend []jsfItem

func (b byPublishedDescend) Len() int {
	return len(b)
}

func (b byPublishedDescend) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byPublishedDescend) Less(i, j int) bool {
	ti, _ := time.Parse(time.RFC3339, b[i].DatePublished)
	tj, _ := time.Parse(time.RFC3339, b[j].DatePublished)
	return ti.After(tj)
}

func findArticlePaths(blogPath string) ([]string, error) {
	blogDir, err := ioutil.ReadDir(blogPath)
	if err != nil {
		return nil, err
	}
	itemPaths := make([]string, 0)
	for _, folder := range blogDir {
		curFolderPath := filepath.Join(blogPath, folder.Name())
		curFolderItemPath := filepath.Join(curFolderPath, itemFile)
		if _, err := os.Stat(curFolderItemPath); err == nil {
			itemPaths = append(itemPaths, curFolderPath)
		}
	}
	return itemPaths, nil
}

func channeledProcessArticle(tmpl *template.Template, articlePath string, ch chan<- jsfItemErr) {
	item, err := processArticle(tmpl, articlePath, "", "")
	ch <- jsfItemErr{item, err}
}

func buildItemList(tmpl *template.Template, blogPath string) ([]jsfItem, error) {
	articlePaths, err := findArticlePaths(blogPath)
	if err != nil {
		return nil, err
	}
	itemList := make([]jsfItem, len(articlePaths))
	ch := make(chan jsfItemErr)
	for _, articlePath := range articlePaths {
		go channeledProcessArticle(tmpl, articlePath, ch)
	}
	for i := range itemList {
		res := <-ch
		if res.err != nil {
			return nil, res.err
		}
		itemList[i] = res.item
	}
	return itemList, nil
}

func (jf *jsfMain) init() error {
	jf.Version = jsfVersion
	jf.Title = hostRawURL
	jf.HomePageURL = hostRawURL

	hostURL, err := url.Parse(hostRawURL)
	if err != nil {
		return err
	}

	URLRelativeToHost, err := url.Parse(jsfPath)
	if err != nil {
		return err
	}

	jf.FeedURL = hostURL.ResolveReference(URLRelativeToHost).String()
	return nil
}

func pageSplit(itemList []jsfItem, pageLen int) ([]jsfMain, error) {
	itemCount := len(itemList)
	feedCount := ((itemCount - 1) / pageLen) + 1
	res := make([]jsfMain, feedCount)
	for i := range res {
		err := res[i].init()
		if err != nil {
			return res, err
		}
		pageStart := i * pageLen
		pageEnd := (i + 1) * pageLen
		if pageEnd > itemCount {
			pageEnd = itemCount
		}
		res[i].Items = itemList[pageStart:pageEnd]
		if i < (feedCount - 1) {
			res[i].NextURL = res[i].FeedURL + strconv.Itoa(i+1)
		}
	}
	return res, nil
}

func writeJsf(feedList []jsfMain, blogPath string) error {
	for i, feed := range feedList {
		curPath := filepath.Join(blogPath, jsfPath)
		if i > 0 {
			curPath += strconv.Itoa(i)
		}
		f, err := os.Create(curPath)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(f)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "\t")
		err = enc.Encode(feed)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}
