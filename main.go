package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
)

// Token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 123456

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimeters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// Lexer

type Lexer struct {
	input        string
	position     int  // Current position in input (points to current char)
	readPosition int  // Current reading position in input (after current char)
	ch           byte // Current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch ch := l.ch; {
	case isLetter(ch):
		tok.Literal = l.readIdentifier()
		tok.Type = LookupIdent(tok.Literal)
		return tok
	case isDigit(ch):
		tok.Type = INT
		tok.Literal = l.readNumber()
	default:
		l.readChar()
		switch ch {
		case '=':
			if l.ch == '=' {
				l.readChar()
				tok = Token{Type: EQ, Literal: "=="}
			} else {
				tok = newToken(ASSIGN, ch)
			}
		case '+':
			tok = newToken(PLUS, ch)
		case '-':
			tok = newToken(MINUS, ch)
		case '!':
			if l.ch == '=' {
				l.readChar()
				tok = Token{Type: NOT_EQ, Literal: "!="}
			} else {
				tok = newToken(BANG, ch)
			}
		case '/':
			tok = newToken(SLASH, ch)
		case '*':
			tok = newToken(ASTERISK, ch)
		case '<':
			tok = newToken(LT, ch)
		case '>':
			tok = newToken(GT, ch)
		case ';':
			tok = newToken(SEMICOLON, ch)
		case ',':
			tok = newToken(COMMA, ch)
		case '(':
			tok = newToken(LPAREN, ch)
		case ')':
			tok = newToken(RPAREN, ch)
		case '{':
			tok = newToken(LBRACE, ch)
		case '}':
			tok = newToken(RBRACE, ch)
		case 0:
			tok.Literal = ""
			tok.Type = EOF
		default:
			tok = newToken(ILLEGAL, ch)
		}
	}

	return tok
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' && ch == '_'
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// Repl

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := NewLexer(line)

		for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}

// Entrypoint

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	Start(os.Stdin, os.Stdout)
}
