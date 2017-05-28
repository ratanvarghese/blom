package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Specify a mode")
	}

	ta, fTemplate := makeTemplateArgs()

	switch os.Args[1] {
	case templateMode:
		if err := fTemplate.Parse(os.Args[2:]); err == nil {
			runTemplate(ta)
		}
	default:
		log.Fatal("Unsupported mode")
	}
}
