package main

import (
	"sort"
	"testing"
	"time"
)

func TestSortByPublished(t *testing.T) {
	t1, _ := time.Parse("2006-01-02", "2015-06-09")
	t2, _ := time.Parse("2006-01-02", "2016-06-09")
	t3, _ := time.Parse("2006-01-02", "2016-07-09")
	t4, _ := time.Parse("2006-01-02", "2016-07-10")

	var ji1, ji2, ji3, ji4, ji5 jsfItem
	ji1.DatePublished = t1.Format(time.RFC3339)
	ji2.DatePublished = t2.Format(time.RFC3339)
	ji3.DatePublished = t3.Format(time.RFC3339)
	ji4.DatePublished = t4.Format(time.RFC3339)
	//ji5.DatePublished stays at zero value

	inOrder := []jsfItem{ji4, ji3, ji2, ji1, ji5}
	scrambled := []jsfItem{ji3, ji1, ji4, ji5, ji2}
	sort.Sort(byPublishedDescend(scrambled))

	for i, expectedItem := range inOrder {
		actualItem := scrambled[i]
		expected := expectedItem.DatePublished
		actual := actualItem.DatePublished
		if expected != actual {
			t.Errorf("Unexpected value at index %v, expected '%s', actual '%s'", i, expected, actual)
		}
	}
}
