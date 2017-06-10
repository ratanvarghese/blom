package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
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

func templateToWriter(wr io.Writer, published time.Time, title, templateText, content string) error {
	return errors.New("Not yet implemented")
}
