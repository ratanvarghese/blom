package main

import (
	"flag"
	"html/template"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Specify a mode")
	}

	aa, fArticle := makeArticleArgs()

	fXarticle := flag.NewFlagSet(xarticleMode, flag.ContinueOnError)
	templateSrc := fXarticle.String("template", defaultTemplate, "Filename of template file")
	tagList := fXarticle.String("tags", "", "Comma-seperated list of tags")
	title := fXarticle.String("title", "", "Title of the article")
	articlePath := fXarticle.String("articledir", ".", "Directory holding the article")

	switch os.Args[1] {
	case articleMode:
		if err := fArticle.Parse(os.Args[2:]); err == nil {
			buildArticle(aa)
		} else {
			log.Fatal(err.Error())
		}
	case updateMode:
		doUpdate()
	case xarticleMode:
		if err := fXarticle.Parse(os.Args[2:]); err == nil {
			tmpl, err := template.ParseFiles(*templateSrc)
			if err != nil {
				log.Fatal(err.Error())
			}

			_, err = processArticle(tmpl, *articlePath, *title, *tagList)
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
	default:
		log.Fatal("Unsupported mode")
	}
}
