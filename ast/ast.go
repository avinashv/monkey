package ast

import "monkey/token"

// Node represents a node in the AST.
type Node interface {
	TokenLiteral() string
	String() string
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

func (program *Program) String() string {
	var output string

	for _, statement := range program.Statements {
		output += statement.String()
	}

	return output
}

// ExpressionStatement represents an expression statement in the AST.
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (expressionStatement *ExpressionStatement) statementNode() {}
func (expressionStatement *ExpressionStatement) TokenLiteral() string {
	return expressionStatement.Token.Literal
}

func (expressionStatement *ExpressionStatement) String() string {
	if expressionStatement.Expression != nil {
		return expressionStatement.Expression.String()
	}

	return ""
}

// Identifier represents an identifier in the AST.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (identifier *Identifier) String() string { return identifier.Value }

func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }

// LetStatement represents a let statement in the AST.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (letStatement *LetStatement) String() string {
	var output string

	output += letStatement.TokenLiteral() + " "
	output += letStatement.Name.String()
	output += " = "

	if letStatement.Value != nil {
		output += letStatement.Value.String()
	}

	output += ";"

	return output
}

func (letStatement *LetStatement) statementNode()       {}
func (letStatement *LetStatement) TokenLiteral() string { return letStatement.Token.Literal }

// ReturnStatement represents a return statement in the AST.
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (returnStatement *ReturnStatement) String() string {
	var output string

	output += returnStatement.TokenLiteral() + " "

	if returnStatement.ReturnValue != nil {
		output += returnStatement.ReturnValue.String()
	}

	output += ";"

	return output
}

func (returnStatement *ReturnStatement) statementNode()       {}
func (returnStatement *ReturnStatement) TokenLiteral() string { return returnStatement.Token.Literal }
