package main

import (
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"strings"
	"time"
)

type articleMeta struct {
	src          string
	url          string
	pagetitle    string
	articletitle string
	publish      bool
	date         string
	topics       []string
}

func main() {
	t := time.Now()
	long := tqtime.LongDate(t.Year(), t.YearDay())
	longish := strings.Replace(long, "After Tranquility", "AT", 1)
	fmt.Printf("%s\n", longish)

	gregorian := t.Format("Monday, 2 January, 2006 CE")
	fmt.Printf("%s\n", gregorian)
}
