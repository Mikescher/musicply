package main

import (
	"git.blackforestbytes.com/BlackForestBytes/goext/bfcodegen"
	"os"
)

func main() {
	dest := os.Args[2]

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = bfcodegen.GenerateCharsetIDSpecs(wd, dest, bfcodegen.CSIDGenOptions{})
	if err != nil {
		panic(err)
	}
}
