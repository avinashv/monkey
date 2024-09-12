package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// Define the precedence of the operators.
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// Define the prefix and infix parse functions.
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser represents the parser.
type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// registerPrefix registers a prefix parse function for a token type.
func (parser *Parser) registerPrefix(tokenType token.TokenType, function prefixParseFn) {
	parser.prefixParseFns[tokenType] = function
}

// registerInfix registers an infix parse function for a token type.
func (parser *Parser) registerInfix(tokenType token.TokenType, function infixParseFn) {
	parser.infixParseFns[tokenType] = function
}

// New creates a new parser instance.
func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer:  lexer,
		errors: []string{},
	}

	parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	parser.registerPrefix(token.IDENT, parser.parseIdentifier)

	// read two tokens, so currentToken and peekToken are both set
	parser.nextToken()
	parser.nextToken()

	return parser
}

// Errors returns the list of errors encountered during parsing.
func (parser *Parser) Errors() []string {
	return parser.errors
}

// peekError appends an error message to the list of errors.
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
		return parser.parseExpressionStatement()
	}
}

// parseExpression parses an expression.
func (parser *Parser) parseExpression(precedence int) ast.Expression {
	// get the prefix parse function for the current token
	prefix := parser.prefixParseFns[parser.currentToken.Type]
	if prefix == nil {
		return nil
	}

	// parse the left expression
	left := prefix()

	// loop until the precedence of the next token is less than the current precedence
	//for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
	//	// get the infix parse function for the next token
	//	infix := parser.infixParseFn[parser.peekToken.Type]
	//	if infix == nil {
	//		return left
	//	}
	//
	//	// advance the tokens
	//	parser.nextToken()
	//
	//	// parse the infix expression
	//	left = infix(left)
	//}

	return left
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

// parseExpressionStatement parses an expression statement.
func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// create the expression statement
	statement := &ast.ExpressionStatement{Token: parser.currentToken}

	// parse the expression
	statement.Expression = parser.parseExpression(LOWEST)

	// check if the next token is a semicolon
	if parser.peekTokenIs(token.SEMICOLON) {
		parser.nextToken()
	}

	// return the expression statement
	return statement
}

// parseIdentifier parses an identifier.
func (parser *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}
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
