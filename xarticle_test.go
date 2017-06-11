package main

import (
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"os"
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
