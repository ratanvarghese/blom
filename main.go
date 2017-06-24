package main

import (
	"flag"
	"html/template"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Specify a mode: '%s' or '%s'", articleMode, updateMode)
	}
	fArticle := flag.NewFlagSet(articleMode, flag.ContinueOnError)
	templateSrc := fArticle.String("template", "../../template.html", "Filename of template file")
	tagList := fArticle.String("tags", "", "Comma-seperated list of tags")
	title := fArticle.String("title", "", "Title of the article")
	articlePath := fArticle.String("articledir", ".", "Directory holding the article")

	fUpdate := flag.NewFlagSet(updateMode, flag.ContinueOnError)
	mainTemplateSrc := fUpdate.String("mtemplate", "../template.html", "Filename of main template file")
	homeTemplateSrc := fUpdate.String("htemplate", "../home-template.html", "Filename of homepage template file")
	blogPath := fUpdate.String("blogdir", ".", "Directory holding the blog")

	switch os.Args[1] {
	case articleMode:
		if err := fArticle.Parse(os.Args[2:]); err == nil {
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
	case updateMode:
		if err := fUpdate.Parse(os.Args[2:]); err == nil {
			mainTmpl, err := template.ParseFiles(*mainTemplateSrc)
			if err != nil {
				log.Fatal(err.Error())
			}
			homeTmpl, err := template.ParseFiles(*homeTemplateSrc)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = processBlog(mainTmpl, homeTmpl, *blogPath)
			if err != nil {
				log.Fatal(err.Error())
			}
		} else {
			log.Fatal(err.Error())
		}
	default:
		log.Fatalf("Unsupported mode: use '%s' or '%s'", articleMode, updateMode)
	}
}
