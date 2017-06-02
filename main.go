package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Specify a mode")
	}

	aa, fArticle := makeArticleArgs()

	switch os.Args[1] {
	case articleMode:
		if err := fArticle.Parse(os.Args[2:]); err == nil {
			buildArticle(aa)
		}
	case updateMode:
		doUpdate()
	default:
		log.Fatal("Unsupported mode")
	}
}
