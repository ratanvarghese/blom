package main

import (
	"fmt"
	"time"
	"strings"
	"github.com/ratanvarghese/tqtime"
)

func main() {
	t := time.Now()
	long := tqtime.LongDate(t.Year(), t.YearDay())
	longish := strings.Replace(long, "After Tranquility", "AT", 1)
	fmt.Printf("%s\n", longish)

	gregorian := t.Format("Monday, 2 January, 2006 CE")
	fmt.Printf("%s\n", gregorian)
}
