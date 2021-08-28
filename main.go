package main

import (
	"./parser"
	"os"
)

func main () {
	vm := parser.InitVm(os.Args[1])
	vm.Run()
}

