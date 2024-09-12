package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">>> "

// Start initializes the REPL.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		// read input from the user
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()

		// check if the user has entered any input or exits the REPL
		if !scanned || scanner.Text() == "" || scanner.Text() == "exit" {
			return
		}

		// lex the input
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		// print the tokens
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

// printParserErrors prints the parser errors to the output.
func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
