package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/ratanvarghese/tqtime"
	"html/template"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const updateMode = "update"
const blogTitle = "ratan.blog"
const jsfVersion = "https://jsonfeed.org/version/1"
const jsfPath = "feeds/json"
const atomPath = "feeds/atom"
const rssPath = "feeds/rss"

type jsfMain struct {
	Version     string    `json:"version"`
	Title       string    `json:"title"`
	HomePageURL string    `json:"home_page_url"`
	FeedURL     string    `json:"feed_url"`
	NextURL     string    `json:"next_url,omitempty"`
	Items       []jsfItem `json:"items"`
}

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

func (jf *jsfMain) init() error {
	jf.Version = jsfVersion
	jf.Title = blogTitle
	jf.HomePageURL = hostRawURL

	hostURL, err := url.Parse(hostRawURL)
	if err != nil {
		return err
	}

	URLRelativeToHost, err := url.Parse(jsfPath)
	if err != nil {
		return err
	}

	jf.FeedURL = hostURL.ResolveReference(URLRelativeToHost).String()
	return nil
}

func pageSplit(itemList []jsfItem, pageLen int) ([]jsfMain, error) {
	itemCount := len(itemList)
	feedCount := ((itemCount - 1) / pageLen) + 1
	res := make([]jsfMain, feedCount)
	for i := range res {
		err := res[i].init()
		if err != nil {
			return res, err
		}
		pageStart := i * pageLen
		pageEnd := (i + 1) * pageLen
		if pageEnd > itemCount {
			pageEnd = itemCount
		}
		res[i].Items = itemList[pageStart:pageEnd]
		if i < (feedCount - 1) {
			res[i].NextURL = res[i].FeedURL + strconv.Itoa(i+1)
		}
	}
	return res, nil
}

func writeJsf(feedList []jsfMain, blogPath string) error {
	for i, feed := range feedList {
		curPath := filepath.Join(blogPath, jsfPath)
		if i > 0 {
			curPath += strconv.Itoa(i)
		}
		f, err := os.Create(curPath)
		if err != nil {
			return err
		}
		enc := json.NewEncoder(f)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "\t")
		err = enc.Encode(feed)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func processHomepage(tmpl *template.Template, wg *sync.WaitGroup, latest jsfItem, blogPath string, ch chan<- error) {
	var exportArgs articleExport
	published, _ := time.Parse(time.RFC3339, latest.DatePublished)
	permalink := fmt.Sprintf("<br /><a href=\"%s\">[Permalink]</a>", latest.URL)
	exportArgs.init(published, latest.Title, []byte(latest.ContentHTML+permalink))
	err := exportArgs.writeFinalWebpage(tmpl, blogPath)
	if err != nil {
		ch <- err
	}
	wg.Done()
}

func archiveSeperator(gt1 time.Time, gt2 time.Time) (bool, string) {
	g1Year := gt1.Year()
	g1YearDay := gt1.YearDay()
	tq1Year := tqtime.Year(g1Year, g1YearDay)
	tq1Mon := tqtime.Month(g1Year, g1YearDay)
	tq1Day := tqtime.Day(g1Year, g1YearDay)

	g2Year := gt2.Year()
	g2YearDay := gt2.YearDay()
	tq2Year := tqtime.Year(g2Year, g2YearDay)
	tq2Mon := tqtime.Month(g2Year, g2YearDay)
	tq2Day := tqtime.Day(g2Year, g2YearDay)

	isSpecialDay := (tq2Mon == tqtime.SpecialDay)

	var seperatorText string
	if tq2Day == tqtime.AldrinDay || tq2Mon == tqtime.Hippocrates { //A feature of this calendar which is annoying for archives
		seperatorText = fmt.Sprintf("<h3>Hippocrates & Aldrin Day, %d AT</h3>", tq2Year)
		if (tq1Mon == tqtime.Hippocrates || tq1Day == tqtime.AldrinDay) && tq1Year == tq2Year {
			return false, seperatorText
		}
	} else if isSpecialDay {
		seperatorText = fmt.Sprintf("<h3>%s, %d AT</h3>", tqtime.DayName(tq2Day), tq2Year)
	} else {
		seperatorText = fmt.Sprintf("<h3>%s, %d AT</h3>", tq2Mon.String(), tq2Year)
	}
	needSeperation := (tq1Year != tq2Year) || (tq1Mon != tq2Mon) || (isSpecialDay && (tq1Day != tq2Day))
	return needSeperation, seperatorText
}

func archiveLines(itemList []jsfItem) []string {
	if len(itemList) < 1 {
		return nil
	}
	var t1 time.Time //intentionally starting at zero value, always a different year than first article.
	outputLines := make([]string, 0)
	for i, ji := range itemList {
		t2, _ := time.Parse(time.RFC3339, ji.DatePublished)
		if sep, sepText := archiveSeperator(t1, t2); sep {
			if i > 0 { //The start of a section is the end of the previous section, unless *no* previous section.
				outputLines = append(outputLines, "</ul>")
			}
			outputLines = append(outputLines, sepText)
			outputLines = append(outputLines, "<ul>")
		}
		outputLines = append(outputLines, fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", ji.URL, ji.Title))
		t1 = t2
	}
	outputLines = append(outputLines, "</ul>")
	return outputLines
}

func processArchive(tmpl *template.Template, wg *sync.WaitGroup, itemList []jsfItem, blogPath string, ch chan<- error) {
	var exportArgs articleExport
	var published time.Time

	contentLines := archiveLines(itemList)
	exportArgs.init(published, "Archive", []byte(strings.Join(contentLines, "\n")))
	exportArgs.Date = template.HTML("")
	archivePath := filepath.Join(blogPath, "archive")
	err := exportArgs.writeFinalWebpage(tmpl, archivePath)
	if err != nil {
		ch <- err
	}
	wg.Done()
}

func tagSort(itemList []jsfItem) (map[string][]jsfItem, []string) {
	res := make(map[string][]jsfItem)
	tagList := make([]string, 0)
	for _, ji := range itemList {
		for _, tag := range ji.Tags {
			if len(res[tag]) == 0 && len(tag) > 0 {
				tagList = append(tagList, tag)
			}
			res[tag] = append(res[tag], ji)
		}
	}
	sort.Strings(tagList)
	return res, tagList
}

func tagsPageLines(itemList []jsfItem) []string {
	outputLines := make([]string, 0)
	tagMap, tagList := tagSort(itemList)
	for _, tag := range tagList {
		outputLines = append(outputLines, fmt.Sprintf("<h3>%v</h3>", strings.Title(tag)))
		outputLines = append(outputLines, "<ul>")
		for _, ji := range tagMap[tag] {
			outputLines = append(outputLines, fmt.Sprintf("<li><a href=\"%v\">%v</a></li>", ji.URL, ji.Title))
		}
		outputLines = append(outputLines, "</ul>")
	}
	return outputLines
}

func processTags(tmpl *template.Template, wg *sync.WaitGroup, itemList []jsfItem, blogPath string, ch chan<- error) {
	var exportArgs articleExport
	var published time.Time

	contentLines := tagsPageLines(itemList)
	exportArgs.init(published, "Tags", []byte(strings.Join(contentLines, "\n")))
	exportArgs.Date = template.HTML("")
	tagsPath := filepath.Join(blogPath, "tags")
	err := exportArgs.writeFinalWebpage(tmpl, tagsPath)
	if err != nil {
		ch <- err
	}
	wg.Done()
}

func fromJsfItem(gi *feeds.Item, ji jsfItem) {
	gi.Title = ji.Title
	gi.Link = &feeds.Link{Href: ji.URL}
	gi.Created, _ = time.Parse(time.RFC3339, ji.DatePublished)
	gi.Updated, _ = time.Parse(time.RFC3339, ji.DateModified)
	gi.Id = ji.URL
	gi.Description = ji.ContentHTML
}

func makeLegacyFeed(itemList []jsfItem) feeds.Feed {
	var gf feeds.Feed
	gf.Title = blogTitle
	gf.Link = &feeds.Link{Href: hostRawURL}
	gf.Created = time.Now()

	gfItemList := make([]*feeds.Item, len(itemList))
	for i, ji := range itemList {
		gfItemList[i] = new(feeds.Item)
		fromJsfItem(gfItemList[i], ji)
	}
	gf.Items = gfItemList
	return gf
}

func processLegacyFeeds(wg *sync.WaitGroup, itemList []jsfItem, blogPath string, ch chan<- error) {
	defer wg.Done()
	gf := makeLegacyFeed(itemList)
	atom, err := gf.ToAtom()
	if err != nil {
		ch <- err
		return
	}
	rss, err := gf.ToRss()
	if err != nil {
		ch <- err
		return
	}

	fullAtomPath := filepath.Join(blogPath, atomPath)
	fullRssPath := filepath.Join(blogPath, atomPath)

	err = ioutil.WriteFile(fullAtomPath, []byte(atom), 0664)
	if err != nil {
		ch <- err
		return
	}
	err = ioutil.WriteFile(fullRssPath, []byte(rss), 0664)
	if err != nil {
		ch <- err
	}
}

func processJsf(wg *sync.WaitGroup, itemList []jsfItem, blogPath string, pageLen int, ch chan<- error) {
	defer wg.Done()
	feedList, err := pageSplit(itemList, pageLen)
	if err != nil {
		ch <- err
		return
	}
	err = writeJsf(feedList, blogPath)
	if err != nil {
		ch <- err
		return
	}
}

func processBlog(mainTmpl *template.Template, homeTmpl *template.Template, blogRelativePath string) error {
	blogPath, err := filepath.Abs(blogRelativePath)
	if err != nil {
		return err
	}

	itemList, err := buildItemList(mainTmpl, blogPath)
	sort.Sort(byPublishedDescend(itemList))

	ch := make(chan error)
	var wg sync.WaitGroup
	if len(itemList) > 0 {
		wg.Add(1)
		go processHomepage(homeTmpl, &wg, itemList[0], blogPath, ch)
	}
	wg.Add(4)
	go processLegacyFeeds(&wg, itemList, blogPath, ch)
	go processTags(mainTmpl, &wg, itemList, blogPath, ch)
	go processArchive(mainTmpl, &wg, itemList, blogPath, ch)
	go processJsf(&wg, itemList, blogPath, 15, ch)
	wg.Wait()
	select {
	case err = <-ch:
		return err
	default:
		return nil
	}
}
