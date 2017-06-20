package main

import (
	"io/ioutil"
	"path/filepath"
	"sort"
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

func TestFindArticlePaths(t *testing.T) {
	numDirs := 3
	numItems := 2
	fileContent := "Fake!"
	blogPath, subdirPaths := setupBlog(t, []byte(fileContent), []byte(fileContent), numDirs, numItems)
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
