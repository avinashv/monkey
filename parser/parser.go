package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
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

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

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
	parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)

	parser.infixParseFns = make(map[token.TokenType]infixParseFn)
	parser.registerInfix(token.PLUS, parser.parseInfixExpression)
	parser.registerInfix(token.MINUS, parser.parseInfixExpression)
	parser.registerInfix(token.SLASH, parser.parseInfixExpression)
	parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfix(token.EQ, parser.parseInfixExpression)
	parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfix(token.LT, parser.parseInfixExpression)
	parser.registerInfix(token.GT, parser.parseInfixExpression)

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
		parser.noPrefixParseFnError(parser.currentToken.Type)
		return nil
	}

	// parse the left expression
	left := prefix()

	// loop until the precedence of the next token is less than the current precedence
	for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
		// get the infix parse function for the next token
		infix := parser.infixParseFns[parser.peekToken.Type]
		if infix == nil {
			return left
		}

		// advance the tokens
		parser.nextToken()

		// parse the infix expression
		left = infix(left)
	}

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

// parseIntegerLiteral parses an integer literal.
func (parser *Parser) parseIntegerLiteral() ast.Expression {
	// create the integer literal
	literal := &ast.IntegerLiteral{Token: parser.currentToken}

	// parse the integer value
	value, err := strconv.ParseInt(parser.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", parser.currentToken.Literal)
		parser.errors = append(parser.errors, msg)
		return nil
	}
	literal.Value = value

	// return the integer literal
	return literal
}

// parsePrefixExpression parses a prefix expression.
func (parser *Parser) parsePrefixExpression() ast.Expression {
	// create the prefix expression
	expression := &ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}

	// advance the tokens
	parser.nextToken()

	// parse the right expression
	expression.Right = parser.parseExpression(PREFIX)

	// return the prefix expression
	return expression
}

// parseInfixExpression parses an infix expression.
func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// create the infix expression
	expression := &ast.InfixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
		Left:     left,
	}

	// get the precedence of the current token
	precedence := parser.currentPrecedence()

	// advance the tokens
	parser.nextToken()

	// parse the right expression
	expression.Right = parser.parseExpression(precedence)

	// return the infix expression
	return expression
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

// peekPrecedence returns the precedence of the peek token.
func (parser *Parser) peekPrecedence() int {
	if precedence, ok := precedences[parser.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

// currentPrecedence returns the precedence of the current token.
func (parser *Parser) currentPrecedence() int {
	if precedence, ok := precedences[parser.currentToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

// noPrefixParseFnError appends an error message to the list of errors.
func (parser *Parser) noPrefixParseFnError(tokenType token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	parser.errors = append(parser.errors, msg)
}
