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

func (identifier *Identifier) String() string       { return identifier.Value }
func (identifier *Identifier) expressionNode()      {}
func (identifier *Identifier) TokenLiteral() string { return identifier.Token.Literal }

// IntegerLiteral represents an integer literal in the AST.
type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (integerLiteral *IntegerLiteral) String() string       { return integerLiteral.Token.Literal }
func (integerLiteral *IntegerLiteral) expressionNode()      {}
func (integerLiteral *IntegerLiteral) TokenLiteral() string { return integerLiteral.Token.Literal }

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

// PrefixExpression represents a prefix expression in the AST.
type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (prefixExpression *PrefixExpression) String() string {
	var output string

	output = "(" + prefixExpression.Operator
	output += prefixExpression.Right.String()
	output += ")"

	return output
}

func (prefixExpression *PrefixExpression) expressionNode() {}
func (prefixExpression *PrefixExpression) TokenLiteral() string {
	return prefixExpression.Token.Literal
}

// InfixExpression represents an infix expression in the AST.
type InfixExpression struct {
	Token    token.Token // the operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (infixExpression *InfixExpression) String() string {
	var output string

	output = "("
	output += infixExpression.Left.String()
	output += " " + infixExpression.Operator + " "
	output += infixExpression.Right.String()
	output += ")"

	return output
}

func (infixExpression *InfixExpression) expressionNode() {}
func (infixExpression *InfixExpression) TokenLiteral() string {
	return infixExpression.Token.Literal
}
