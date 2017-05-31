package main

import (
	"fmt"
	"io/ioutil"
)

const updateMode = "update"

func foo() {
	dir, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}
	for _, file := range dir {
		fmt.Println(file.Name())
	}
}
