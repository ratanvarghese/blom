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
	templateSrc := fXarticle.String("template", "../../template.html", "Filename of template file")
	tagList := fXarticle.String("tags", "", "Comma-seperated list of tags")
	title := fXarticle.String("title", "", "Title of the article")
	articlePath := fXarticle.String("articledir", ".", "Directory holding the article")

	fXupdate := flag.NewFlagSet(xupdateMode, flag.ContinueOnError)
	mainTemplateSrc := fXupdate.String("mtemplate", "../template.html", "Filename of main template file")
	homeTemplateSrc := fXupdate.String("htemplate", "../home-template.html", "Filename of homepage template file")
	blogPath := fXupdate.String("blogdir", ".", "Directory holding the blog")

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
	case xupdateMode:
		if err := fXupdate.Parse(os.Args[2:]); err == nil {
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
		log.Fatal("Unsupported mode")
	}
}
