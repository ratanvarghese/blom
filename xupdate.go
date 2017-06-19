package main

import (
	"time"
)

type byDatePublished []jsfItem

func (b byDatePublished) Len() int {
	return len(b)
}

func (b byDatePublished) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byDatePublished) Less(i, j int) bool {
	ti, _ := time.Parse(time.RFC3339, b[i].DatePublished)
	tj, _ := time.Parse(time.RFC3339, b[j].DatePublished)
	return ti.After(tj)
}
