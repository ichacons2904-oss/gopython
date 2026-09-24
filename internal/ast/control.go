package ast

import "gopython/internal/lexer"

type IfStatement struct {
	Token     lexer.Token
	Condition Expression
	Body      []Statement
	Elif      []ElifBranch
	Else      []Statement
}

func (node IfStatement) Position() lexer.Token {
	return node.Token
}

func (IfStatement) statementNode() {}

type ElifBranch struct {
	Token     lexer.Token
	Condition Expression
	Body      []Statement
}

func (node ElifBranch) Position() lexer.Token {
	return node.Token
}

type WhileStatement struct {
	Token     lexer.Token
	Condition Expression
	Body      []Statement
}

func (node WhileStatement) Position() lexer.Token {
	return node.Token
}

func (WhileStatement) statementNode() {}

type ForStatement struct {
	Token    lexer.Token
	Target   Identifier
	Iterable Expression
	Body     []Statement
}

func (node ForStatement) Position() lexer.Token {
	return node.Token
}

func (ForStatement) statementNode() {}

type ReturnStatement struct {
	Token lexer.Token
	Value Expression
}

func (node ReturnStatement) Position() lexer.Token {
	return node.Token
}

func (ReturnStatement) statementNode() {}

type BreakStatement struct {
	Token lexer.Token
}

func (node BreakStatement) Position() lexer.Token {
	return node.Token
}

func (BreakStatement) statementNode() {}

type ContinueStatement struct {
	Token lexer.Token
}

func (node ContinueStatement) Position() lexer.Token {
	return node.Token
}

func (ContinueStatement) statementNode() {}
