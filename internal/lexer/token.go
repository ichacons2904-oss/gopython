package lexer

type TokenType string

const (
	EOF TokenType = "EOF"

	Identifier TokenType = "IDENTIFIER"
	Integer    TokenType = "INTEGER"
	String     TokenType = "STRING"

	If       TokenType = "IF"
	Else     TokenType = "ELSE"
	Elif     TokenType = "ELIF"
	While    TokenType = "WHILE"
	For      TokenType = "FOR"
	In       TokenType = "IN"
	Def      TokenType = "DEF"
	Return   TokenType = "RETURN"
	Break    TokenType = "BREAK"
	Continue TokenType = "CONTINUE"
	And      TokenType = "AND"
	Or       TokenType = "OR"
	Not      TokenType = "NOT"
	True     TokenType = "TRUE"
	False    TokenType = "FALSE"
	None     TokenType = "NONE"

	Newline TokenType = "NEWLINE"
	Indent  TokenType = "INDENT"
	Dedent  TokenType = "DEDENT"

	Assign       TokenType = "ASSIGN"
	Plus         TokenType = "PLUS"
	Minus        TokenType = "MINUS"
	Asterisk     TokenType = "ASTERISK"
	Slash        TokenType = "SLASH"
	Percent      TokenType = "PERCENT"
	Equal        TokenType = "EQUAL"
	NotEqual     TokenType = "NOT_EQUAL"
	LessThan     TokenType = "LESS_THAN"
	LessEqual    TokenType = "LESS_EQUAL"
	GreaterThan  TokenType = "GREATER_THAN"
	GreaterEqual TokenType = "GREATER_EQUAL"

	Colon            TokenType = "COLON"
	LeftParenthesis  TokenType = "LEFT_PAREN"
	RightParenthesis TokenType = "RIGHT_PAREN"
	Comma            TokenType = "COMMA"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
	Column int
}

var keywords = map[string]TokenType{
	"if":       If,
	"else":     Else,
	"elif":     Elif,
	"while":    While,
	"for":      For,
	"in":       In,
	"def":      Def,
	"return":   Return,
	"break":    Break,
	"continue": Continue,
	"and":      And,
	"or":       Or,
	"not":      Not,
	"True":     True,
	"False":    False,
	"None":     None,
}

var twoCharacterTokens = map[string]TokenType{
	"==": Equal,
	"!=": NotEqual,
	"<=": LessEqual,
	">=": GreaterEqual,
}

var oneCharacterTokens = map[rune]TokenType{
	'=': Assign,
	'+': Plus,
	'-': Minus,
	'*': Asterisk,
	'/': Slash,
	'%': Percent,
	'<': LessThan,
	'>': GreaterThan,
	':': Colon,
	'(': LeftParenthesis,
	')': RightParenthesis,
	',': Comma,
}

func tokenTypeForIdentifier(lexeme string) TokenType {
	if tokenType, ok := keywords[lexeme]; ok {
		return tokenType
	}
	return Identifier
}
