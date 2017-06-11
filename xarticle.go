package main

import (
	"errors"
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type jsfAttachment struct {
	URL      string `json:"url"`
	MIMEType string `json:"mime_type"`
	valid    bool
}

type jsfItem struct {
	ID            string          `json:"id"`
	URL           string          `json:"url"`
	Title         string          `json:"title"`
	ContentHTML   string          `json:"content_html"`
	DatePublished string          `json:"date_published"`
	DateModified  string          `json:"date_modified"`
	Tags          []string        `json:"tags"`
	Attachments   []jsfAttachment `json:"attachments"`
}

type articleExport struct {
	Title       string
	Date        template.HTML
	Today       template.HTML
	ContentHTML template.HTML
}

const hostRawURL = "http://ratan.blog"
const attachmentDir = "attachments"
const listSeperator = ","

func (ja *jsfAttachment) init(basename string, article string, fileStart []byte) error {
	ja.MIMEType = http.DetectContentType(fileStart)
	URLRelativeToHost, err := url.Parse(article + "/" + attachmentDir + "/" + basename)
	if err != nil {
		return err
	}

	hostURL, err := url.Parse(hostRawURL)
	if err != nil {
		return err
	}

	ja.URL = hostURL.ResolveReference(URLRelativeToHost).String()

	ja.valid = true
	return nil
}

func (ji *jsfItem) init(published, modified time.Time, title, directory, tagList string) error {
	if len(title) < 1 {
		return errors.New("Blank title")
	}

	if len(directory) < 1 {
		return errors.New("Blank directory")
	}

	base, err := url.Parse(hostRawURL)
	if err != nil {
		return err
	}

	u, err := url.Parse(filepath.Base(directory))
	if err != nil {
		return err
	}

	ji.URL = base.ResolveReference(u).String()
	ji.ID = ji.URL
	ji.Title = title
	ji.DatePublished = published.Format(time.RFC3339)
	if published.After(modified) {
		ji.DateModified = ji.DatePublished
	} else {
		ji.DateModified = modified.Format(time.RFC3339)
	}
	if len(tagList) > 0 {
		ji.Tags = strings.Split(tagList, listSeperator)
	}
	return nil
}

func (articleE *articleExport) init(published time.Time, title string, content []byte) {
	articleE.Title = title
	articleE.Date = template.HTML(dualDateStr(published))
	articleE.Today = "Today is " + template.HTML(dualDateStr(time.Now()))
	articleE.ContentHTML = template.HTML(content)
}

func dualDateStr(gDate time.Time) string {
	const outputGDateFormat = "Monday, 2 January, 2006 CE"
	tqDate := tqtime.LongDate(gDate.Year(), gDate.YearDay())
	tqDateBetter := strings.Replace(tqDate, "After Tranquility", "AT", 1)
	gDateStr := gDate.Format(outputGDateFormat)
	return fmt.Sprintf("%s<br />[Gregorian: %s]", tqDateBetter, gDateStr)
}

func getAttachPaths(articlePath string) (map[string]bool, error) {
	attachPath := filepath.Join(articlePath, attachmentDir)
	attachList, err := ioutil.ReadDir(attachPath)
	if err != nil {
		return nil, err
	}

	res := make(map[string]bool)
	for _, attachFileInfo := range attachList {
		if !attachFileInfo.IsDir() {
			res[filepath.Join(attachPath, attachFileInfo.Name())] = true
		}
	}

	return res, nil
}

func attachmentsFromReaders(article string, filepaths []string, buffers []io.Reader) ([]jsfAttachment, error) {
	if len(filepaths) != len(buffers) {
		return nil, errors.New("Mismatch between filepath count and buffer count.")
	}

	const bytesNeededToFindMIMEType = 512
	attachList := make([]jsfAttachment, len(filepaths))
	for i, curpath := range filepaths {
		b := make([]byte, bytesNeededToFindMIMEType)
		_, err := buffers[i].Read(b)
		if err != nil {
			return attachList, err
		}
		attachList[i].init(filepath.Base(curpath), article, b)
	}

	return attachList, nil
}
