package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var jpegBytes []byte

func TestMain(m *testing.M) {
	jpegBytes = []byte{255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 255, 219, 0, 67, 0, 6, 4, 5, 6, 5, 4, 6, 6, 5, 6, 7, 7, 6, 8, 10, 16, 10, 10, 9, 9, 10, 20, 14, 15, 12, 16, 23, 20, 24, 24, 23, 20, 22, 22, 26, 29, 37, 31, 26, 27, 35, 28, 22, 22, 32, 44, 32, 35, 38, 39, 41, 42, 41, 25, 31, 45, 48, 45, 40, 48, 37, 40, 41, 40, 255, 219, 0, 67, 1, 7, 7, 7, 10, 8, 10, 19, 10, 10, 19, 40, 26, 22, 26, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 255, 192, 0, 17, 8, 2, 163, 4, 176, 3, 1, 34, 0, 2, 17, 1, 3, 17, 1, 255, 196, 0, 31, 0, 0, 1, 5, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 255, 196, 0, 181, 16, 0, 2, 1, 3, 3, 2, 4, 3, 5, 5, 4, 4, 0, 0, 1, 125, 1, 2, 3, 0, 4, 17, 5, 18, 33, 49, 65, 6, 19, 81, 97, 7, 34, 113, 20, 50, 129, 145, 161, 8, 35, 66, 177, 193, 21, 82, 209, 240, 36, 51, 98, 114, 130, 9, 10, 22, 23, 24, 25, 26, 37, 38, 39, 40, 41, 42, 52, 53, 54, 55, 56, 57, 58, 67, 68, 69, 70, 71, 72, 73, 74, 83, 84, 85, 86, 87, 88, 89, 90, 99, 100, 101, 102, 103, 104, 105, 106, 115, 116, 117, 118, 119, 120, 121, 122, 131, 132, 133, 134, 135, 136, 137, 138, 146, 147, 148, 149, 150, 151, 152, 153, 154, 162, 163, 164, 165, 166, 167, 168, 169, 170, 178, 179, 180, 181, 182, 183, 184, 185, 186, 194, 195, 196, 197, 198, 199, 200, 201, 202, 210, 211, 212, 213, 214, 215, 216, 217, 218, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 255, 196, 0, 31, 1, 0, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 255, 196, 0, 181, 17, 0, 2, 1, 2, 4, 4, 3, 4, 7, 5, 4, 4, 0, 1, 2, 119, 0, 1, 2, 3, 17, 4, 5, 33, 49, 6, 18, 65, 81, 7, 97, 113, 19, 34, 50, 129, 8, 20, 66, 145, 161, 177, 193, 9, 35, 51, 82, 240, 21, 98, 114, 209, 10, 22, 36, 52, 225, 37, 241, 23, 24, 25, 26, 38, 39, 40, 41, 42, 53, 54, 55, 56, 57, 58, 67, 68, 69, 70, 71, 72, 73}
	os.Exit(m.Run())
}

func TestAttachInitJPEG(t *testing.T) {
	baseName := "1200.jpg"
	articleName := "hello"
	var ja jsfAttachment
	err := ja.init(baseName, articleName, jpegBytes)

	if err != nil {
		t.Errorf("Error (%s) when given valid inputs.", err.Error())
	}

	expectedMIME := "image/jpeg"
	if ja.MIMEType != expectedMIME {
		t.Errorf("Wrong MIME Type, expected:%s, actual:%s", expectedMIME, ja.MIMEType)
	}

	expectedURL := fmt.Sprintf("%s/%s/attachments/%s", hostRawURL, articleName, baseName)
	if ja.URL != expectedURL {
		t.Errorf("Wrong URL, expected:%s, actual:%s", expectedURL, ja.URL)
	}

	if !ja.valid {
		t.Errorf("jsfAttachment is invalid after given valid arguments")
	}
}

func TestItemInitPubAfterMod(t *testing.T) {
	published, err := time.Parse("2006-01-02", "2017-06-09")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	modified, err := time.Parse("2006-01-02", "2017-06-08")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	var ji jsfItem
	err = ji.init(published, modified, "Hey", "hey", "")
	if err != nil {
		t.Errorf("Error (%s) for valid input.", err.Error())
	}

	expectedPub := published.Format(time.RFC3339)
	if ji.DatePublished != expectedPub {
		t.Errorf("Wrong publish date, expected '%s', actual '%s'.", expectedPub, ji.DatePublished)
	}

	if ji.DateModified != expectedPub {
		t.Errorf("Not converting modified date to publish date, when publish date is later.")
	}
}

func TestItemInitBlanks(t *testing.T) {
	published, err := time.Parse("2006-01-02", "2017-06-09")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	modified, err := time.Parse("2006-01-02", "2017-06-08")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	var ji1 jsfItem
	err = ji1.init(published, modified, "", "hey", "")
	if err == nil {
		t.Errorf("No error for blank title")
	}

	var ji2 jsfItem
	err = ji2.init(published, modified, "Hey", "", "")
	if err == nil {
		t.Errorf("No error for blank directory")
	}

	var ji3 jsfItem
	err = ji3.init(published, modified, "Hey", "hey", "")
	if err != nil {
		t.Errorf("Error (%s) for blank tags", err.Error())
	}

	if len(ji3.Tags) > 0 {
		t.Errorf("%v tags when none provided.", len(ji3.Tags))
	}

}

func TestItemInit(t *testing.T) {
	published, err := time.Parse("2006-01-02", "2017-06-08")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	modified, err := time.Parse("2006-01-02", "2017-06-09")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	title := "demo title"
	directory := "demo"
	tagList := "demo,test,xxx"

	var ji jsfItem
	err = ji.init(published, modified, title, directory, tagList)
	if err != nil {
		t.Errorf("Error (%s) for valid input.", err.Error())
	}

	expectedPub := published.Format(time.RFC3339)
	if ji.DatePublished != expectedPub {
		t.Errorf("Wrong publish date, expected '%s', actual '%s'.", expectedPub, ji.DatePublished)
	}

	expectedMod := modified.Format(time.RFC3339)
	if ji.DateModified != expectedMod {
		t.Errorf("Wrong modification date, expected '%s', actual '%s'.", expectedMod, ji.DateModified)
	}

	if ji.Title != title {
		t.Errorf("Wrong title, expected '%s', actual '%s'.", title, ji.Title)
	}

	expectedURL := fmt.Sprintf("%s/%s", hostRawURL, directory)
	if ji.URL != expectedURL {
		t.Errorf("Wrong URL, expected '%s', actual '%s'.", expectedURL, ji.URL)
	}
	if ji.ID != expectedURL {
		t.Errorf("Wrong ID, expected '%s', actual '%s'.", expectedURL, ji.ID)
	}

	expectedTagCount := 3
	actualTagCount := len(ji.Tags)
	if actualTagCount < expectedTagCount {
		t.Errorf("Wrong number of tags, expected %v, actual %v", expectedTagCount, actualTagCount)
	} else if ji.Tags[0] != "demo" || ji.Tags[1] != "test" || ji.Tags[2] != "xxx" {
		t.Errorf("Wrong tags")
	}
}

func TestArticleExportInit(t *testing.T) {
	title := "Demo Title"
	miniContent := "Demo Content"

	published, err := time.Parse("2006-01-02", "2017-06-10")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	var articleE articleExport
	articleE.init(published, title, []byte(miniContent))

	if articleE.Title != title {
		t.Errorf("Wrong title, expected '%s', actual '%s'.", title, articleE.Title)
	}

	today := time.Now()
	gregString := today.Format("Monday, 2 January, 2006 CE")
	tqString := tqtime.LongDate(today.Year(), today.YearDay())
	tqStringBetter := strings.Replace(tqString, "After Tranquility", "AT", 1)
	expectedTodayString := fmt.Sprintf("Today is %s<br />[Gregorian: %s]", tqStringBetter, gregString)
	if string(articleE.Today) != expectedTodayString {
		t.Errorf("Wrong 'today' string, expected '%s', actual '%s'.", expectedTodayString, articleE.Today)
	}
	expectedDate := "Sunday, 17 Lavoisier, 48 AT<br />[Gregorian: Saturday, 10 June, 2017 CE]"
	if string(articleE.Date) != expectedDate {
		t.Errorf("Wrong 'date' string, expected '%s', actual '%s'.", expectedDate, articleE.Date)
	}

	if string(articleE.ContentHTML) != miniContent {
		t.Errorf("Wrong content, expected '%s', actual '%s'.", miniContent, articleE.ContentHTML)
	}

}

func TestDualDateStr(t *testing.T) {
	input, err := time.Parse("2006-01-02", "2017-06-10")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	result := dualDateStr(input)
	expected := "Sunday, 17 Lavoisier, 48 AT<br />[Gregorian: Saturday, 10 June, 2017 CE]"
	if result != expected {
		t.Errorf("Wrong date, expected '%s', actual '%s'.", expected, result)
	}
}

func setupArticlePath(t *testing.T) string {
	articlePath, err := ioutil.TempDir(".", "testblom")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	return articlePath
}

func setupAttachPaths(t *testing.T) (string, string, map[string]bool) {
	articlePath := setupArticlePath(t)
	attachPath := filepath.Join(articlePath, attachmentDir)
	err := os.Mkdir(attachPath, 0777)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	jpeg1Path := filepath.Join(attachPath, "jpeg_1.jpeg")
	err = ioutil.WriteFile(jpeg1Path, jpegBytes, 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	jpeg2Path := filepath.Join(attachPath, "jpeg_2.jpeg")
	err = ioutil.WriteFile(jpeg2Path, jpegBytes, 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	ignoreMe := filepath.Join(articlePath, "ignoreMe.jpeg")
	err = ioutil.WriteFile(ignoreMe, jpegBytes, 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	expectedAttachPaths := make(map[string]bool)
	expectedAttachPaths[jpeg1Path] = true
	expectedAttachPaths[jpeg2Path] = true

	return articlePath, attachPath, expectedAttachPaths
}

func teardownArticlePath(t *testing.T, articlePath string) {
	err := os.RemoveAll(articlePath)
	if err != nil {
		t.Errorf("Error (%s) AFTER RUNNING TEST", err.Error())
	}

}

func TestGetAttachPaths(t *testing.T) {
	articlePath, _, expectedAttachPaths := setupAttachPaths(t)
	attachPaths, err := getAttachPaths(articlePath)
	if err != nil {
		t.Errorf("Error (%s) for valid inputs.", err.Error())
	}

	expectedPathCount := len(expectedAttachPaths)
	actualPathCount := len(attachPaths)
	if actualPathCount != expectedPathCount {
		t.Errorf("Wrong number of attachment paths, expected %v, actual %v.", expectedPathCount, actualPathCount)
	}

	for someAttachPath := range attachPaths {
		if !expectedAttachPaths[someAttachPath] {
			t.Errorf("Unexpected path '%s'.", someAttachPath)
		}
	}

	for someAttachPath := range expectedAttachPaths {
		if !attachPaths[someAttachPath] {
			t.Errorf("Missing path '%s'.", someAttachPath)
		}
	}
	teardownArticlePath(t, articlePath)
}

func TestAttachmentsFromReaders(t *testing.T) {
	article := "demo"
	filenames := []string{"a1.jpeg", "a2.jpeg", "a3.jpeg"}
	buf1 := bytes.NewBuffer(jpegBytes)
	buf2 := bytes.NewBuffer(jpegBytes)
	buf3 := bytes.NewBuffer(jpegBytes)
	buffers := []io.Reader{buf1, buf2, buf3}

	attachments, err := attachmentsFromReaders(article, filenames, buffers)
	if err != nil {
		t.Errorf("Error (%s) with valid inputs.", err.Error())
	}

	for i, filename := range filenames {
		var expected jsfAttachment
		err = expected.init(filename, article, jpegBytes)
		if err != nil {
			t.Errorf("Error (%s) with attachment init.", err.Error())
		}
		if attachments[i] != expected {
			t.Errorf("[At index %v] expected %v, actual %v.", i, expected, attachments[i])
		}
	}
}

func TestGetArticleContentMD(t *testing.T) {
	articlePath := setupArticlePath(t)
	MDContent := "## This is a heading"
	MDPath := filepath.Join(articlePath, contentFileMD)
	err := ioutil.WriteFile(MDPath, []byte(MDContent), 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	HTMLContent := "Ignore Me"
	HTMLPath := filepath.Join(articlePath, contentFileHTML)
	err = ioutil.WriteFile(HTMLPath, []byte(HTMLContent), 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	articleContent, _, err := getArticleContent(articlePath)
	if err != nil {
		t.Errorf("Error (%s) with valid inputs.", err.Error())
	}

	expectedArticleContent := "<h2>This is a heading</h2>\n"
	if string(articleContent) != expectedArticleContent {
		t.Errorf("Wrong content, expected '%s', actual '%s'.", expectedArticleContent, string(articleContent))
	}
	teardownArticlePath(t, articlePath)
}

func TestGetArticleContentHTML(t *testing.T) {
	articlePath := setupArticlePath(t)
	HTMLContent := "Do NOT Ignore Me"
	HTMLPath := filepath.Join(articlePath, contentFileHTML)
	err := ioutil.WriteFile(HTMLPath, []byte(HTMLContent), 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	articleContent, _, err := getArticleContent(articlePath)
	if err != nil {
		t.Errorf("Error (%s) with valid inputs.", err.Error())
	}

	if string(articleContent) != HTMLContent {
		t.Errorf("Wrong content, expected '%s', actual '%s'.", HTMLContent, string(articleContent))
	}
	teardownArticlePath(t, articlePath)
}

func TestGetPreviousItemNonexistent(t *testing.T) {
	articlePath := setupArticlePath(t)
	_, fileExists, err := getPreviousItem(articlePath)
	if err != nil {
		t.Errorf("Error (%s) with valid inputs.", err.Error())
	}
	if fileExists {
		t.Errorf("fileExists == true when item file does not exist.")
	}
	teardownArticlePath(t, articlePath)
}

func TestGetPreviousItemExistent(t *testing.T) {
	articlePath := setupArticlePath(t)
	published, err := time.Parse("2006-01-02", "2017-06-08")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	modified, err := time.Parse("2006-01-02", "2017-06-09")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	var ji jsfItem
	err = ji.init(published, modified, "demo title", "demo", "demo,test,xxx")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	itemFilePath := filepath.Join(articlePath, itemFile)
	f, err := os.Create(itemFilePath)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	err = enc.Encode(ji)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	jiOutput, fileExists, err := getPreviousItem(articlePath)
	if err != nil {
		t.Errorf("Error (%s) with valid inputs.", err.Error())
	}
	if !fileExists {
		t.Errorf("fileExists != true when item file exists.")
	}
	if jiOutput.ID != ji.ID {
		t.Errorf("Wrong ID, expected '%s', actual '%s'.", ji.ID, jiOutput.ID)
	}
	if jiOutput.Title != ji.Title {
		t.Errorf("Wrong title, expected '%s', actual '%s'.", ji.Title, jiOutput.Title)
	}
	if jiOutput.DateModified != ji.DateModified {
		t.Errorf("Wrong date modified, expected '%v', actual '%v',", ji.DateModified, jiOutput.DateModified)
	}
	if jiOutput.DatePublished != ji.DatePublished {
		t.Errorf("Wrong date published, expected '%v', actual '%v',", ji.DatePublished, jiOutput.DatePublished)
	}

	teardownArticlePath(t, articlePath)
}

func TestFilesFromAttachPathMap(t *testing.T) {
	articlePath := setupArticlePath(t)
	jpegFile1 := filepath.Join(articlePath, "file1.jpeg")
	jpegFile2 := filepath.Join(articlePath, "file2.jpeg")
	err := ioutil.WriteFile(jpegFile1, jpegBytes, 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	err = ioutil.WriteFile(jpegFile2, jpegBytes, 0664)
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}

	attachPathMap := make(map[string]bool)
	attachPathMap[jpegFile1] = true
	attachPathMap[jpegFile2] = true
	expectedLen := len(attachPathMap)

	attachPathList, attachFileList, attachReaderList, err := filesFromAttachPathMap(attachPathMap)
	if err != nil {
		t.Errorf("Error (%s) when all parameters valid")
	}

	attachPathListLen := len(attachPathList)
	if attachPathListLen != expectedLen {
		t.Errorf("Wrong attachPathList length, expected %v, actual %v", expectedLen, attachPathListLen)
	}
	for path := range attachPathMap {
		if attachPathList[0] != path && attachPathList[1] != path {
			t.Errorf("Missing path '%v' in attachPathList", path)
		}
	}
	for i, path := range attachPathList {
		if !attachPathMap[path] {
			t.Errorf("Unexpected path '%v' at index %v in attachPathList", path, i)
		}
	}

	attachFileListLen := len(attachFileList)
	if attachFileListLen != expectedLen {
		t.Errorf("Wrong attachFileList length, expected %v, actual %v", expectedLen, attachFileListLen)
	}
	for path := range attachPathMap {
		if attachFileList[0].Name() != path && attachFileList[1].Name() != path {
			t.Errorf("Missing path '%v' in attachPathList", path)
		}
	}
	for i, file := range attachFileList {
		if !attachPathMap[file.Name()] {
			t.Errorf("Unexpected path '%v' at index %v in attachFileList", file.Name(), i)
		}
	}

	attachReaderListLen := len(attachReaderList)
	if attachReaderListLen != expectedLen {
		t.Errorf("Wrong attachReaderList length, expected %v, actual %v", expectedLen, attachReaderListLen)
	}

	teardownArticlePath(t, articlePath)
}

func TestWriteItemFile(t *testing.T) {
	var ji jsfItem
	ji.ID = "Test ID"
	ji.URL = "http://example.com"
	ji.Title = "Test Title"
	ji.ContentHTML = "<h1>Hello World!</h1>"
	ji.DatePublished = time.Now().Format(time.RFC3339)
	ji.DateModified = time.Now().Format(time.RFC3339)
	ji.Tags = []string{"tag1", "testtag"}

	articlePath := setupArticlePath(t)
	err := writeItemFile(ji, articlePath)
	if err != nil {
		t.Errorf("Error (%s) when all parameters valid.", err.Error())
	}

	itemFilePath := filepath.Join(articlePath, itemFile)
	itemFileContent, err := ioutil.ReadFile(itemFilePath)
	if err != nil {
		t.Errorf("Error (%s) when reading item file", err.Error())
	}

	var res jsfItem
	json.Unmarshal(itemFileContent, &res)
	if res.ID != ji.ID {
		t.Errorf("Wrong ID, expected '%s', actual '%s'", ji.ID, res.ID)
	}
	if res.URL != ji.URL {
		t.Errorf("Wrong URL, expected '%s', actual '%s'", ji.URL, res.URL)
	}
	if res.Title != ji.Title {
		t.Errorf("Wrong Title, expected '%s', actual '%s'", ji.Title, res.Title)
	}
	if res.ContentHTML != ji.ContentHTML {
		t.Errorf("Wrong ContentHTML, expected '%s', actual '%s'", ji.ContentHTML, res.ContentHTML)
	}
	if res.DatePublished != ji.DatePublished {
		t.Errorf("Wrong DatePublished, expected '%s', actual '%s'", ji.DatePublished, res.DatePublished)
	}
	if res.DateModified != ji.DateModified {
		t.Errorf("Wrong DateModified, expected '%s', actual '%s'", ji.DateModified, res.DateModified)
	}
	if res.Tags[0] != ji.Tags[0] || res.Tags[1] != ji.Tags[1] {
		t.Errorf("Wrong tags, expected '%v', actual '%v'", ji.Tags, res.Tags)
	}

	teardownArticlePath(t, articlePath)
}
