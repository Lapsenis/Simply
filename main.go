package main

import (
	"Simply/interpreter"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) > 1 {
		interpreter.ProcessFile(os.Args[1], os.Stdout)
	} else {
		fmt.Println("Simply 0.1")
		interpreter.Start(os.Stdin, os.Stdout)
	}
}
