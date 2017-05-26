package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type ArticleExport struct {
	Title        string
	Content_html template.HTML
}

func demoTemplate() {
	t, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	content, err := ioutil.ReadFile("hello.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	article1 := ArticleItem{}
	article1.Id = "http://ratan.blog/hello"
	article1.Url = "http://ratan.blog/hello"
	article1.Title = "Hello"
	article1.Content_html = string(content)
	article1.Date_published = "2017-05-25T8:04:00-05:00"
	article1.Date_modified = "2017-05-25T8:11:00-11:30"
	article1.Tags = []string{"nonsense", "meta"}

	f, err := os.Create("output.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	w := bufio.NewWriter(f)

	articleE := ArticleExport{article1.Title, template.HTML(article1.Content_html)}

	err = t.Execute(w, articleE)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Flush()
}

func demoDate() {
	t := time.Now()
	long := tqtime.LongDate(t.Year(), t.YearDay())
	longish := strings.Replace(long, "After Tranquility", "AT", 1)
	fmt.Println(longish)

	gregorian := t.Format("Monday, 2 January, 2006 CE")
	fmt.Println(gregorian)
}

func demoFeed() {
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
		return
	}
	fmt.Println(string(b))
}
