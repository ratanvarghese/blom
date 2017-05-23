package main

import (
	"encoding/json"
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"io/ioutil"
	"strings"
	"time"
)

type ArticleMeta struct {
	Src          string
	URL          string
	PageTitle    string
	ArticleTitle string
	Publish      bool
	Date         string
	Topics       []string
}

func main() {
	t := time.Now()
	long := tqtime.LongDate(t.Year(), t.YearDay())
	longish := strings.Replace(long, "After Tranquility", "AT", 1)
	fmt.Println(longish)

	gregorian := t.Format("Monday, 2 January, 2006 CE")
	fmt.Println(gregorian)

	filename := "data.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	var aml []ArticleMeta
	err = json.Unmarshal(data, &aml)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, element := range aml {
		fmt.Println(element.ArticleTitle)
	}
}
