package repl

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
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

		// print the tokens
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
