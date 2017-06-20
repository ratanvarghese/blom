package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

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
