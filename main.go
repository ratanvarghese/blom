package main

import (
	"flag"
)

func main() {
	ta := makeTemplateArgs()

	flag.Parse()

	runTemplate(ta)
}
