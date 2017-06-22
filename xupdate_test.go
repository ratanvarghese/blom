package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

func TestJsfMainInit(t *testing.T) {
	var jf jsfMain
	err := jf.init()
	if err != nil {
		t.Errorf("Error (%s) with default settings.", err.Error())
	}
	if jf.Version != jsfVersion {
		t.Errorf("Wrong version, expected '%s', actual '%s'", jsfVersion, jf.Version)
	}
	if jf.Title != hostRawURL {
		t.Errorf("Wrong title, expected '%s', actual '%s'", hostRawURL, jf.Title)
	}
	if jf.HomePageURL != hostRawURL {
		t.Errorf("Wrong home URL, expected '%s', actual '%s'", hostRawURL, jf.HomePageURL)
	}
}

var pageSplitTestParams = []struct {
	itemCount int
	pageLen   int
}{
	{1, 3},
	{2, 3},
	{3, 3},
	{5, 3},
	{6, 3},
	{7, 3},
	{4, 15},
	{15, 15},
	{16, 15},
	{60, 15},
	{61, 15},
}

func pageSplitTest(t *testing.T, itemCount int, pageLen int) {
	itemList := make([]jsfItem, itemCount)

	for i := range itemList {
		itemList[i].ID = strconv.Itoa(i)
	}
	feedList, err := pageSplit(itemList, pageLen)

	if err != nil {
		t.Errorf("Error (%s) when all parameters valid.", err.Error())
		t.Errorf("(itemCount %v, pageLen %v)", itemCount, pageLen)
	}

	expectedFeedCount := ((itemCount - 1) / pageLen) + 1
	feedCount := len(feedList)
	if feedCount != expectedFeedCount {
		t.Errorf("Wrong feed count, expected %v, actual %v", expectedFeedCount, feedCount)
		t.Errorf("(itemCount %v, pageLen %v)", itemCount, pageLen)
	}

	for fi, feed := range feedList {
		for i, item := range feed.Items {
			expectedID := strconv.Itoa(fi*pageLen + i)
			ID := item.ID
			if ID != expectedID {
				t.Errorf("Wrong ID in feed %v, index %v, expected '%s', actual '%s'", fi, i, expectedID, ID)
				t.Errorf("(itemCount %v, pageLen %v)", itemCount, pageLen)
			}
		}
		if feedCount > 1 && fi < (feedCount-1) {
			expectedNext := feed.FeedURL + strconv.Itoa(fi+1)
			if feed.NextURL != expectedNext {
				t.Errorf("Wrong NextURL in feed %v, expected '%s', actual '%s'", fi, expectedNext, feed.NextURL)
			}
		} else if feedCount > 1 && fi == feedCount-1 && len(feed.NextURL) > 0 {
			t.Errorf("NextURL in feed %v when not expected", fi)
		}
	}
}

func TestPageSplitAll(t *testing.T) {
	for _, s := range pageSplitTestParams {
		pageSplitTest(t, s.itemCount, s.pageLen)
	}
}

func TestWriteJsf(t *testing.T) {
	feedCount := 3
	feedList := make([]jsfMain, feedCount)
	for i := range feedList {
		feedList[i].Title = strconv.Itoa(i)
	}

	blogPath, _ := setupBlog(t, []byte("fake"), []byte("fake"), 0, 0)

	err := os.Mkdir(filepath.Join(blogPath, "feeds"), 0777)
	if err != nil {
		t.Errorf("Error (%s) BEFORE RUNNING TEST", err.Error())
	}
	err = writeJsf(feedList, blogPath)
	if err != nil {
		t.Errorf("Error (%s) when all parameters valid.", err.Error())
	}

	feedFileList := make([]string, feedCount)
	for i := range feedFileList {
		if i > 0 {
			feedFileList[i] = filepath.Join(blogPath, jsfPath+strconv.Itoa(i))
		} else {
			feedFileList[i] = filepath.Join(blogPath, jsfPath)
		}
	}

	for i, feedFileName := range feedFileList {
		if _, err = os.Stat(feedFileName); err != nil {
			t.Errorf("Error (%s) seeking file '%s'", err.Error(), feedFileName)
		}
		feedBytes, err := ioutil.ReadFile(feedFileName)
		if err != nil {
			t.Errorf("Error (%s) accessing file '%s'", err.Error(), feedFileName)
		}
		var curFeed jsfMain
		err = json.Unmarshal(feedBytes, &curFeed)
		if err != nil {
			t.Errorf("Error (%s) decoding file '%s'", err.Error(), feedFileName)
		}
		if curFeed.Title != feedList[i].Title {
			t.Errorf("Wrong title at index %v, expected '%s', actual '%s'", i, feedList[i].Title, curFeed.Title)
		}
	}
	teardownArticlePath(t, blogPath)
}

var archiveSeperatorTests = []struct {
	gt1     string
	gt2     string
	sep     bool
	sepText string
}{
	{"0001-01-01", "1972-07-20", true, "<h3>Armstrong Day, 3 AT</h3>"},
	{"1972-07-19", "1972-07-20", true, "<h3>Armstrong Day, 3 AT</h3>"},
	{"0001-01-01", "1972-06-22", true, "<h3>Mendel, 3 AT</h3>"},
	{"1972-06-21", "1972-06-22", true, "<h3>Mendel, 3 AT</h3>"},
	{"1972-06-22", "1972-06-23", false, "<h3>Mendel, 3 AT</h3>"},
	{"1972-01-04", "1972-01-05", true, "<h3>Galileo, 3 AT</h3>"},
	{"1972-02-01", "1972-02-02", true, "<h3>Hippocrates & Aldrin Day, 3 AT</h3>"},
	{"1972-02-28", "1972-02-29", false, "<h3>Hippocrates & Aldrin Day, 3 AT</h3>"},
	{"1971-02-28", "1972-02-29", true, "<h3>Hippocrates & Aldrin Day, 3 AT</h3>"},
	{"1972-02-29", "1972-03-01", false, "<h3>Hippocrates & Aldrin Day, 3 AT</h3>"},
	{"1972-03-01", "1972-03-02", true, "<h3>Imhotep, 3 AT</h3>"},
}

func TestArchiveSeperator(t *testing.T) {
	for _, s := range archiveSeperatorTests {
		gt1, _ := time.Parse("2006-01-02", s.gt1)
		gt2, _ := time.Parse("2006-01-02", s.gt2)
		sep, sepText := archiveSeperator(gt1, gt2)
		if sep != s.sep {
			t.Errorf("Wrong seperation status on (%v,%v), expected %v, actual %v", s.gt1, s.gt2, s.sep, sep)
		}
		if sepText != s.sepText {
			t.Errorf("Wrong seperation text on (%v,%v), expected %v, actual %v", s.gt1, s.gt2, s.sepText, sepText)
		}
	}
}

func TestArchiveLinesEmpty(t *testing.T) {
	lineList := archiveLines(nil)
	if len(lineList) > 0 {
		t.Errorf("Got %v with nil input.", lineList)
	}
}

func TestArchiveLinesNonempty(t *testing.T) {
	t0, _ := time.Parse("2006-01-02", "1972-03-01")
	t1, _ := time.Parse("2006-01-02", "1972-03-02")
	t2, _ := time.Parse("2006-01-02", "1972-03-03")
	itemList := make([]jsfItem, 3)

	itemList[0].DatePublished = t0.Format(time.RFC3339)
	itemList[0].Title = "item0 Title"
	itemList[0].URL = "http://ratan.blog/0"
	itemList[1].DatePublished = t1.Format(time.RFC3339)
	itemList[1].Title = "item1 Title"
	itemList[1].URL = "http://ratan.blog/1"
	itemList[2].DatePublished = t2.Format(time.RFC3339)
	itemList[2].Title = "item2 Title"
	itemList[2].URL = "http://ratan.blog/2"

	lineList := archiveLines(itemList)

	expectedLineList := []string{
		"<h3>Hippocrates & Aldrin Day, 3 AT</h3>",
		"<ul>",
		fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", itemList[0].URL, itemList[0].Title),
		"</ul>",
		"<h3>Imhotep, 3 AT</h3>",
		"<ul>",
		fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", itemList[1].URL, itemList[1].Title),
		fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", itemList[2].URL, itemList[2].Title),
		"</ul>",
	}

	expectedLen := len(expectedLineList)
	actualLen := len(lineList)
	if expectedLen != actualLen {
		t.Errorf("Unexpected line count, expected %v, actual %v", expectedLen, actualLen)
	}

	for i, line := range lineList {
		expectedLine := expectedLineList[i]
		if line != expectedLine {
			t.Errorf("Unexpected value at index %v, expected '%s', actual '%s'", i, expectedLine, line)
		}
	}
}

func TestTagSort(t *testing.T) {
	itemList := make([]jsfItem, 4)
	itemList[0].Tags = []string{"tag3", "tag2"}
	itemList[0].ID = strconv.Itoa(0)
	itemList[1].Tags = []string{"tag3"}
	itemList[1].ID = strconv.Itoa(1)
	itemList[2].Tags = []string{"tag1"}
	itemList[2].ID = strconv.Itoa(2)
	itemList[3].ID = strconv.Itoa(3)
	//last item deliberately has no tags

	tagMap, tagList := tagSort(itemList)
	expectedTagCount := 3
	mapTagCount := len(tagMap)
	if mapTagCount != expectedTagCount {
		t.Errorf("Wrong tag count for map, expected %v, actual %v", expectedTagCount, mapTagCount)
	}
	listTagCount := len(tagList)
	if listTagCount != expectedTagCount {
		t.Errorf("Wrong tag count for list, expected %v, actual %v", expectedTagCount, listTagCount)
	}

	expectedTagList := []string{"tag1", "tag2", "tag3"}
	expectedID1 := []string{"2"}
	expectedID2 := []string{"0"}
	expectedID3 := []string{"0", "1"}
	expectedIDListList := [][]string{expectedID1, expectedID2, expectedID3}
	for i, tag := range tagList {
		expectedTag := expectedTagList[i]
		if tag != expectedTag {
			t.Errorf("Unexpected tag in list at index %v, expected '%s', actual '%s'", i, expectedTag, tag)
		}
		itemList := tagMap[tag]
		IDList := expectedIDListList[i]
		for j, item := range itemList {
			ID := IDList[j]
			if item.ID != ID {
				t.Errorf("Unexpected id for tag %s, index %v, expected '%s', actual '%s'", tag, j, ID, item.ID)
			}
		}
	}

}

func TestTagsPageLines(t *testing.T) {
	itemList := make([]jsfItem, 4)
	itemList[0].Tags = []string{"tag3", "tag2"}
	itemList[0].URL = strconv.Itoa(0)
	itemList[0].Title = "title0"
	itemList[1].Tags = []string{"tag3"}
	itemList[1].URL = strconv.Itoa(1)
	itemList[1].Title = "title1"
	itemList[2].Tags = []string{"tag1"}
	itemList[2].URL = strconv.Itoa(2)
	itemList[2].Title = "title2"
	itemList[3].URL = strconv.Itoa(3)
	itemList[3].Title = "title3"

	lineList := tagsPageLines(itemList)
	expectedLineList := []string{
		"<h3>Tag1</h3>",
		"<ul>",
		fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", itemList[2].URL, itemList[2].Title),
		"</ul>",
		"<h3>Tag2</h3>",
		"<ul>",
		fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", itemList[0].URL, itemList[0].Title),
		"</ul>",
		"<h3>Tag3</h3>",
		"<ul>",
		fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", itemList[0].URL, itemList[0].Title),
		fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", itemList[1].URL, itemList[1].Title),
		"</ul>",
	}

	lineCount := len(lineList)
	expectedLineCount := len(expectedLineList)
	if lineCount != expectedLineCount {
		t.Errorf("Wrong line count, expected %v, actual %v", expectedLineCount, lineCount)
	}

	for i, line := range lineList {
		expectedLine := expectedLineList[i]
		if line != expectedLine {
			t.Errorf("Unexpected line at index %v, expected '%s', actual '%s", i, expectedLine, line)
		}
	}
}
