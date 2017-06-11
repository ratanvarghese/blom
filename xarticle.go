package main

import (
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"html/template"
	"net/http"
	"net/url"
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
	base, err := url.Parse(hostRawURL)
	if err != nil {
		return err
	}

	u, err := url.Parse(directory)
	if err != nil {
		return err
	}

	ji.URL = base.ResolveReference(u).String()
	ji.ID = ji.URL
	ji.Title = title
	ji.DatePublished = published.Format(time.RFC3339)
	ji.DateModified = modified.Format(time.RFC3339)
	ji.Tags = strings.Split(tagList, listSeperator)
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
