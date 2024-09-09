package main

import (
	"fmt"
	"monkey/repl"
	"os"
)

func main() {
	// initialize the REPL
	fmt.Printf("Monkey v0.1\n")
	repl.Start(os.Stdin, os.Stdout)
}
