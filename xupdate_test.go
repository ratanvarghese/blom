package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestSortByPublished(t *testing.T) {
	t1, _ := time.Parse("2006-01-02", "2015-06-09")
	t2, _ := time.Parse("2006-01-02", "2016-06-09")
	t3, _ := time.Parse("2006-01-02", "2016-07-09")
	t4, _ := time.Parse("2006-01-02", "2016-07-10")

	var ji1, ji2, ji3, ji4, ji5 jsfItem
	ji1.DatePublished = t1.Format(time.RFC3339)
	ji2.DatePublished = t2.Format(time.RFC3339)
	ji3.DatePublished = t3.Format(time.RFC3339)
	ji4.DatePublished = t4.Format(time.RFC3339)
	//ji5.DatePublished stays at zero value

	inOrder := []jsfItem{ji4, ji3, ji2, ji1, ji5}
	scrambled := []jsfItem{ji3, ji1, ji4, ji5, ji2}
	sort.Sort(byPublishedDescend(scrambled))

	for i, expectedItem := range inOrder {
		actualItem := scrambled[i]
		expected := expectedItem.DatePublished
		actual := actualItem.DatePublished
		if expected != actual {
			t.Errorf("Unexpected value at index %v, expected '%s', actual '%s'", i, expected, actual)
		}
	}
}

func setupBlog(t *testing.T, itemFileContent, articleContent []byte, numDirs, numItems int) (string, []string) {
	blogPath, err := ioutil.TempDir(".", "testblom")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	subdirPaths := make([]string, numDirs)
	for i, _ := range subdirPaths {
		subdirPath, err := ioutil.TempDir(blogPath, "testblom")
		if err != nil {
			t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
		}
		subdirPaths[i] = subdirPath
		contentPath := filepath.Join(subdirPath, contentFileMD)
		err = ioutil.WriteFile(contentPath, articleContent, 0664)
		if err != nil {
			t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
		}
		if i < numItems {
			itemPath := filepath.Join(subdirPath, itemFile)
			err = ioutil.WriteFile(itemPath, itemFileContent, 0664)
			if err != nil {
				t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
			}
		}
	}
	return blogPath, subdirPaths
}

func setupArticle(t *testing.T, articlePath string, itemFileContent, articleContent []byte) {
	contentPath := filepath.Join(articlePath, contentFileMD)
	err := ioutil.WriteFile(contentPath, articleContent, 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST", err.Error())
	}
	itemPath := filepath.Join(articlePath, itemFile)
	err = ioutil.WriteFile(itemPath, itemFileContent, 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST", err.Error())
	}
}

func TestFindArticlePaths(t *testing.T) {
	numDirs := 3
	numItems := 2
	blogPath, subdirPaths := setupBlog(t, []byte("Fake!"), []byte("Fake!"), numDirs, numItems)
	expectedArticlePaths := subdirPaths[:numItems]
	articlePaths, err := findArticlePaths(blogPath)
	if err != nil {
		t.Errorf("Error (%s) when all inputs valid.", err.Error())
	}
	expectedArticleCount := len(expectedArticlePaths)
	actualArticleCount := len(articlePaths)
	if actualArticleCount != expectedArticleCount {
		t.Errorf("Wrong number of article paths, expected %v, actual %v", expectedArticleCount, actualArticleCount)
	}

	sort.Strings(expectedArticlePaths)
	sort.Strings(articlePaths)
	for i, actualPath := range articlePaths {
		expectedPath := expectedArticlePaths[i]
		if actualPath != expectedPath {
			t.Errorf("Unexpected path at index %v, expected '%s', actual '%s'", i, expectedPath, actualPath)
		}
	}
	teardownArticlePath(t, blogPath)
}

func TestBuildItemList(t *testing.T) {
	numDirs := 3
	numItems := 2
	fileContent := "Placeholder"
	blogPath, subdirPaths := setupBlog(t, []byte(fileContent), []byte(fileContent), numDirs, numItems)
	articlePaths := subdirPaths[:numItems]
	published0, _ := time.Parse("2006-01-02", "2016-07-09")
	modified0, _ := time.Parse("2006-01-02", "2016-07-10")
	var item0 jsfItem
	item0.init(published0, modified0, "Title 0", filepath.Base(articlePaths[0]), "natural,imperative")
	itemBytes0, _ := json.Marshal(&item0)
	setupArticle(t, articlePaths[0], itemBytes0, []byte("## This is the content 0"))

	published1, _ := time.Parse("2006-01-02", "2016-06-09")
	modified1, _ := time.Parse("2006-01-02", "2016-06-10")
	var item1 jsfItem
	item1.init(published1, modified1, "Title 1", filepath.Base(articlePaths[1]), "natural,imperative")
	itemBytes1, _ := json.Marshal(&item1)
	setupArticle(t, articlePaths[1], itemBytes1, []byte("## This is the content 1"))

	templateStr := "{{.Title}}\n{{.Date}}\n{{.Today}}\n{{.ContentHTML}}"
	tmpl := template.New("Whatever")
	tmpl.Parse(templateStr)

	itemList, err := buildItemList(tmpl, blogPath)
	if err != nil {
		t.Errorf("Error (%s) when all parameters valid.", err.Error())
	}
	sort.Sort(byPublishedDescend(itemList))
	expectedContentList := []string{"<h2>This is the content 0</h2>\n", "<h2>This is the content 1</h2>\n"}
	expectedItemList := []jsfItem{item0, item1}

	if len(itemList) != numItems {
		t.Errorf("Wrong number of items, expected %v, actual %v", len(itemList), numItems)
	}

	for i, actualItem := range itemList {
		actualContent := actualItem.ContentHTML
		expectedContent := expectedContentList[i]
		if actualContent != expectedContent {
			t.Errorf("Wrong content at index %v, expected '%s', actual '%s'", i, expectedContent, actualContent)
		}
		actualPub := actualItem.DatePublished
		expectedPub := expectedItemList[i].DatePublished
		if actualPub != expectedPub {
			t.Errorf("Wrong published date at index $v, expected '%s', actual '%s'", i, expectedPub, actualPub)
		}
		finalPagePath := filepath.Join(articlePaths[i], finalWebpageFile)
		finalPageContent, err := ioutil.ReadFile(finalPagePath)
		if err != nil {
			t.Errorf("Error (%s) finding final page %v", err.Error(), i)
		}
		finalPageLines := strings.Split(string(finalPageContent), "\n")
		if finalPageLines[0] != itemList[i].Title {
			t.Errorf("Wrong title in final file %v, expected '%s', actual '%s'", i, itemList[i].Title, finalPageLines[0])
		}

	}
	teardownArticlePath(t, blogPath)
}
