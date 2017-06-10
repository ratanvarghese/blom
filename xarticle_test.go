package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestAttachInitJPEG(t *testing.T) {
	jpegBytes := []byte{255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 255, 219, 0, 67, 0, 6, 4, 5, 6, 5, 4, 6, 6, 5, 6, 7, 7, 6, 8, 10, 16, 10, 10, 9, 9, 10, 20, 14, 15, 12, 16, 23, 20, 24, 24, 23, 20, 22, 22, 26, 29, 37, 31, 26, 27, 35, 28, 22, 22, 32, 44, 32, 35, 38, 39, 41, 42, 41, 25, 31, 45, 48, 45, 40, 48, 37, 40, 41, 40, 255, 219, 0, 67, 1, 7, 7, 7, 10, 8, 10, 19, 10, 10, 19, 40, 26, 22, 26, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 255, 192, 0, 17, 8, 2, 163, 4, 176, 3, 1, 34, 0, 2, 17, 1, 3, 17, 1, 255, 196, 0, 31, 0, 0, 1, 5, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 255, 196, 0, 181, 16, 0, 2, 1, 3, 3, 2, 4, 3, 5, 5, 4, 4, 0, 0, 1, 125, 1, 2, 3, 0, 4, 17, 5, 18, 33, 49, 65, 6, 19, 81, 97, 7, 34, 113, 20, 50, 129, 145, 161, 8, 35, 66, 177, 193, 21, 82, 209, 240, 36, 51, 98, 114, 130, 9, 10, 22, 23, 24, 25, 26, 37, 38, 39, 40, 41, 42, 52, 53, 54, 55, 56, 57, 58, 67, 68, 69, 70, 71, 72, 73, 74, 83, 84, 85, 86, 87, 88, 89, 90, 99, 100, 101, 102, 103, 104, 105, 106, 115, 116, 117, 118, 119, 120, 121, 122, 131, 132, 133, 134, 135, 136, 137, 138, 146, 147, 148, 149, 150, 151, 152, 153, 154, 162, 163, 164, 165, 166, 167, 168, 169, 170, 178, 179, 180, 181, 182, 183, 184, 185, 186, 194, 195, 196, 197, 198, 199, 200, 201, 202, 210, 211, 212, 213, 214, 215, 216, 217, 218, 225, 226, 227, 228, 229, 230, 231, 232, 233, 234, 241, 242, 243, 244, 245, 246, 247, 248, 249, 250, 255, 196, 0, 31, 1, 0, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 255, 196, 0, 181, 17, 0, 2, 1, 2, 4, 4, 3, 4, 7, 5, 4, 4, 0, 1, 2, 119, 0, 1, 2, 3, 17, 4, 5, 33, 49, 6, 18, 65, 81, 7, 97, 113, 19, 34, 50, 129, 8, 20, 66, 145, 161, 177, 193, 9, 35, 51, 82, 240, 21, 98, 114, 209, 10, 22, 36, 52, 225, 37, 241, 23, 24, 25, 26, 38, 39, 40, 41, 42, 53, 54, 55, 56, 57, 58, 67, 68, 69, 70, 71, 72, 73}
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

func TestTemplateToWriter(t *testing.T) {
	miniTemplate := "{{.Title}}\n{{.Today}}\n{{.Date}}\n{{.ContentHTML}}"
	title := "Demo Title"
	published, err := time.Parse("2006-01-02", "2017-07-07")
	if err != nil {
		t.Errorf("Error (%s) PRIOR TO RUNNING TEST.", err.Error())
	}
	outputBuf := new(bytes.Buffer)
	miniContent := "Demo Content"

	err = templateToWriter(outputBuf, published, title, miniTemplate, miniContent)
	if err != nil {
		t.Errorf("Error (%s) when given valid inputs.", err.Error())
	}

	resultLines := strings.Split(outputBuf.String(), "\n")
	if len(resultLines) < 1 {
		t.Errorf("No title")
	} else if resultLines[0] != title {
		t.Errorf("Wrong title, expected '%s', actual '%s'.", title, resultLines[0])
	}

	if len(resultLines) < 2 {
		t.Errorf("No 'today' string")
	}

	if len(resultLines) < 3 {
		t.Errorf("No 'date' string")
	}

	if len(resultLines) < 4 {
		t.Errorf("No content")
	} else if resultLines[3] != miniContent {
		t.Errorf("Wrong content, expected '%s', actual '%s'.", miniContent, resultLines[3])
	}

}
