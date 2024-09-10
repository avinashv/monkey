package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	lexer        *lexer.Lexer
	errors       []string
	currentToken token.Token
	peekToken    token.Token
}

// New creates a new parser instance.
func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	// read two tokens, so currentToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()

	return parser
}

func (parser *Parser) Errors() []string {
	return parser.errors
}

func (parser *Parser) peekError(token token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", token, parser.peekToken.Type)
	parser.errors = append(parser.errors, msg)
}

// nextToken advances the currentToken and peekToken.
func (parser *Parser) nextToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken = parser.lexer.NextToken()
}

// ParseProgram parses the program.
func (parser *Parser) ParseProgram() *ast.Program {
	// create the root node of the AST
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// parse each statement in the program until EOF token is found
	for parser.currentToken.Type != token.EOF {
		// parse the statement
		statement := parser.parseStatement()

		// add the statement to the program if not nil
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		parser.nextToken()
	}

	// return the program, nothing is left to parse
	return program
}

// parseStatement parses a statement.
func (parser *Parser) parseStatement() ast.Statement {
	switch parser.currentToken.Type {
	case token.LET:
		return parser.parseLetStatement()
	case token.RETURN:
		return parser.parseReturnStatement()
	default:
		return nil
	}
}

// parseLetStatement parses a let statement.
func (parser *Parser) parseLetStatement() *ast.LetStatement {
	// create the let statement
	statement := &ast.LetStatement{Token: parser.currentToken}

	// check if the next token is an identifier
	if !parser.expectPeek(token.IDENT) {
		return nil
	}

	// create the identifier
	statement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	// check if the next token is an assignment
	if !parser.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO skip the expression for now
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	// return the let statement
	return statement
}

// parseReturnStatement parses a return statement.
func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
	// create the return statement
	statement := &ast.ReturnStatement{Token: parser.currentToken}

	// TODO skip the expression for now
	for !parser.currentTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	// return the return statement
	return statement
}

// currentTokenIs checks if the current token is of the given type.
func (parser *Parser) currentTokenIs(tokenType token.TokenType) bool {
	return parser.currentToken.Type == tokenType
}

// peekTokenIs checks if the peek token is of the given type.
func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
	return parser.peekToken.Type == tokenType
}

// expectPeek checks if the peek token is of the given type and advances the tokens.
func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
	if parser.peekTokenIs(tokenType) {
		parser.nextToken()
		return true
	} else {
		parser.peekError(tokenType)
		return false
	}
}
