package main

import "log"

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
