package ast

import "monkey/token"

// Node represents a node in the AST.
type Node interface {
	TokenLiteral() string
}

// Statement represents a statement in the AST.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression in the AST.
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root node of the AST.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the token literal of the first statement in the program.
func (program *Program) TokenLiteral() string {
	if len(program.Statements) > 0 {
		return program.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement represents a let statement in the AST.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }

// Identifier represents an identifier in the AST.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }

// ReturnStatement represents a return statement in the AST.
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (returnStatement *ReturnStatement) statementNode()       {}
func (returnStatement *ReturnStatement) TokenLiteral() string { return returnStatement.Token.Literal }
