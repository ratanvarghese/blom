package main

import (
	"bufio"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	applyTemplate := flag.Bool("apply", false, "Apply a template")
	templateFile := flag.String("template", "template.html", "File to use as a template")
	contentFile := flag.String("content", "content.html", "HTML content to insert in template")
	articleTitle := flag.String("title", "Untitled", "Title of the article")
	outputFile := flag.String("output", "index.html", "Output of the executed template")

	flag.Parse()

	if *applyTemplate {
		runTemplate(*templateFile, *contentFile, *articleTitle, *outputFile)
	} else {
		log.Fatal("Only apply is supported right now!")
	}
}

type ArticleExport struct {
	Title        string
	Content_html template.HTML
}

func runTemplate(templateFile string, contentFile string, articleTitle string, outputFile string) {
	t, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	content, err := ioutil.ReadFile(contentFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	w := bufio.NewWriter(f)

	articleE := ArticleExport{articleTitle, template.HTML(content)}

	err = t.Execute(w, articleE)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Flush()
	f.Close()
}
