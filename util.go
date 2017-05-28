package main

import (
	"fmt"
	"github.com/ratanvarghese/tqtime"
	"log"
	"strings"
	"time"
)

const gDateYYYYMMDDFormat = "2006-01-02"

func killOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printOnError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func todayYYYYMMDD() string {
	return time.Now().Format(gDateYYYYMMDDFormat)
}

func webpageDate(dateYYYYMMDD string) (string, error) {
	const outputGDateFormat = "Monday, 2 January, 2006 CE"
	gDate, err := time.Parse(gDateYYYYMMDDFormat, dateYYYYMMDD)
	if err != nil {
		return "", err
	}

	tqDate := tqtime.LongDate(gDate.Year(), gDate.YearDay())
	tqDateBetter := strings.Replace(tqDate, "After Tranquility", "AT", 1)
	gDateStr := gDate.Format(outputGDateFormat)
	return fmt.Sprintf("%s [Gregorian: %s]", tqDateBetter, gDateStr), err
}

func headerDate() string {
	actualInfo, _ := webpageDate(todayYYYYMMDD())
	return fmt.Sprintf("Today is %s", actualInfo)
}
