package main

import (
	"encoding/json"
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"log"
	"strings"
	"time"
)

type Author struct {
	Name string `json:"name"`
}

type Attachment struct {
	Url       string `json:"url"`
	Mime_type string `json:"mime_type"`
}

type ArticleItem struct {
	Id             string       `json:"id"`
	Url            string       `json:"url"`
	Title          string       `json:"title"`
	Content_html   string       `json:"content_html"`
	Date_published string       `json:"date_published"`
	Date_modified  string       `json:"date_modified"`
	Tags           []string     `json:"tags"`
	Attachments    []Attachment `json:"attachments"`
}

type JSONFeed struct {
	Version       string        `json:"version"`
	Title         string        `json:"title"`
	Home_page_url string        `json:"home_page_url"`
	Feed_url      string        `json:"feed_url"`
	Next_url      string        `json:"next_url"`
	Icon          string        `json:"icon"`
	Favicon       string        `json:"favicon"`
	Author        Author        `json:"author"`
	Items         []ArticleItem `json:"items"`
}

func main() {
	t := time.Now()
	long := tqtime.LongDate(t.Year(), t.YearDay())
	longish := strings.Replace(long, "After Tranquility", "AT", 1)
	fmt.Println(longish)

	gregorian := t.Format("Monday, 2 January, 2006 CE")
	fmt.Println(gregorian)

	j := JSONFeed{}
	j.Version = "https://jsonfeed.org/Version/1"
	j.Title = "ratan.blog"
	j.Home_page_url = "http://ratan.blog"
	j.Feed_url = "http://ratan.blog/feeds/json"
	j.Next_url = "http://ratan.blog/feeds/json1"
	j.Icon = "http://ratan.blog/icon.gif"
	j.Favicon = "http://ratan.blog/favicon.ico"

	a := Author{}
	a.Name = "Ratan Varghese"
	j.Author = a

	article1 := ArticleItem{}
	article1.Id = "http://ratan.blog/hello"
	article1.Url = "http://ratan.blog/hello"
	article1.Title = "Hello"
	article1.Content_html = "<h2>Hello</h2><p>So this is my new blog and stuff</p>"
	article1.Date_published = "2017-05-25T8:04:00-05:00"
	article1.Date_modified = "2017-05-25T8:11:00-11:30"
	article1.Tags = []string{"nonsense", "meta"}

	article2 := ArticleItem{}
	article2.Id = "http://ratan.blog/life-with-a-dumb-phone"
	article2.Url = "http://ratan.blog/life-with-a-dumb-phone"
	article2.Title = "Life With a Dumb Phone"
	article2.Content_html = "<h2>Life With a Dumb Phone</h2><p>This thing is the worst</p>"
	article2.Date_published = "2017-05-26T8:04:00-05:00"
	article2.Date_modified = "2017-05-26T8:11:00-11:30"
	article2.Tags = []string{"observations", "technology"}

	attach2_2 := Attachment{}
	attach2_2.Url = "http://ratan.blog/life-with-a-dumb-phone/lame.jpg"
	attach2_2.Mime_type = "image/jpeg"
	article2.Attachments = []Attachment{attach2_2}

	j.Items = append(j.Items, article1)
	j.Items = append(j.Items, article2)

	b, err := json.MarshalIndent(j, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
