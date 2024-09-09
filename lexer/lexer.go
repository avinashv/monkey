package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

// New creates a new lexer instance.
func New(input string) *Lexer {
	lexer := &Lexer{input: input}

	lexer.readChar()

	return lexer
}

// peekChar returns the next character in the input without advancing the position.
func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		// EOF
		return 0
	} else {
		// peek the next character
		return lexer.input[lexer.readPosition]
	}
}

// readChar reads the next character in the input and advances the position in the input string.
func (lexer *Lexer) readChar() {
	lexer.char = lexer.peekChar()

	// move the position forward
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

// NextToken returns the next token in the input.
func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	// skip whitespace
	lexer.skipWhitespace()

	switch lexer.char {
	case '=':
		// check for equality or assignment
		if lexer.peekChar() == '=' {
			// read the next character
			lexer.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, lexer.char)
		}
	case '+':
		tok = newToken(token.PLUS, lexer.char)
	case '-':
		tok = newToken(token.MINUS, lexer.char)
	case '!':
		// check for inequality or bang
		if lexer.peekChar() == '=' {
			// read the next character
			lexer.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, lexer.char)
		}
	case '/':
		tok = newToken(token.SLASH, lexer.char)
	case '*':
		tok = newToken(token.ASTERISK, lexer.char)
	case '<':
		tok = newToken(token.LT, lexer.char)
	case '>':
		tok = newToken(token.GT, lexer.char)
	case ';':
		tok = newToken(token.SEMICOLON, lexer.char)
	case ',':
		tok = newToken(token.COMMA, lexer.char)
	case '(':
		tok = newToken(token.LPAREN, lexer.char)
	case ')':
		tok = newToken(token.RPAREN, lexer.char)
	case '{':
		tok = newToken(token.LBRACE, lexer.char)
	case '}':
		tok = newToken(token.RBRACE, lexer.char)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(lexer.char) {
			// read the identifier
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lexer.char) {
			// read the number
			tok.Literal = lexer.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			// illegal character
			tok = newToken(token.ILLEGAL, lexer.char)
		}
	}

	lexer.readChar()
	return tok
}

// skipWhitespace skips any whitespace characters in the input.
func (lexer *Lexer) skipWhitespace() {
	for lexer.char == ' ' || lexer.char == '\t' || lexer.char == '\n' || lexer.char == '\r' {
		lexer.readChar()
	}
}

// newToken creates a new token with the given type and character.
func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

// readIdentifier reads an identifier from the input.
func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isLetter(lexer.char) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

// readNumber reads a number from the input.
func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.char) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

// isLetter checks if the given character is a letter.
func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

// isDigit checks if the given character is a digit.
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
