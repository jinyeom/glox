package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type TokenType byte

const (
	// Symbols
	TokLeftParen    TokenType = iota // (
	TokRightParen                    // )
	TokLeftBrace                     // [
	TokRightBrace                    // }
	TokComma                         // ,
	TokDot                           // .
	TokMinus                         // -
	TokPlus                          // +
	TokSemicolon                     // ;
	TokSlash                         // /
	TokStar                          // *
	TokBang                          // !
	TokBangEqual                     // !=
	TokEqual                         // =
	TokEqualEqual                    // ==
	TokGreater                       // >
	TokGreaterEqual                  // >=
	TokLess                          // <
	TokLessEqual                     // <=

	// Literals
	TokIdentifier
	TokString
	TokNumber

	// Reserved keywords
	TokAnd    // and
	TokClass  // class
	TokElse   // else
	TokFalse  // false
	TokFun    // fun
	TokFor    // for
	TokIf     // if
	TokNil    // nil
	TokOr     // or
	TokPrint  // print
	TokReturn // return
	TokSuper  // super
	TokThis   // this
	TokTrue   // true
	TokVar    // var
	TokWhile  // while

	// EOF
	TokEOF
)

// keywords is a map of reserved keywords to their token types
var keywords map[string]TokenType = map[string]TokenType{
	"and":    TokAnd,
	"class":  TokClass,
	"else":   TokElse,
	"false":  TokFalse,
	"for":    TokFor,
	"fun":    TokFun,
	"if":     TokIf,
	"nil":    TokNil,
	"or":     TokOr,
	"print":  TokPrint,
	"return": TokReturn,
	"super":  TokSuper,
	"this":   TokThis,
	"true":   TokTrue,
	"var":    TokVar,
	"while":  TokWhile,
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("[%d]%s:%d(%v)", t.Line, t.Lexeme, t.Type, t.Literal)
}

// ScanTokens scans the argument source code and return scanned tokens.
func ScanTokens(source string) []Token {
	var s scanner.Scanner
	s.Init(strings.NewReader(source))
	tokens := make([]Token, 0)
	for t := s.Scan(); t != scanner.EOF; t = s.Scan() {
		fmt.Println(s.TokenText())
	}
	return tokens
}

// TODO: does this really need to be a structure, or can we get away with a scan function?
type Scanner struct {
	src                  string
	tokens               []Token
	start, current, line int
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		src:     src,
		tokens:  make([]Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]Token, bool) {
	for !s.Done() {
		s.start = s.current
		if ok := s.scanToken(); !ok {
			return nil, ok
		}
	}
	s.tokens = append(s.tokens, Token{TokEOF, "", nil, s.line})
	return s.tokens, true
}

func (s *Scanner) Done() bool {
	return s.current >= len(s.src)
}

func (s *Scanner) scanToken() bool {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(TokLeftParen, nil)
	case ')':
		s.addToken(TokRightParen, nil)
	case '{':
		s.addToken(TokLeftBrace, nil)
	case '}':
		s.addToken(TokRightBrace, nil)
	case ',':
		s.addToken(TokComma, nil)
	case '.':
		s.addToken(TokDot, nil)
	case '-':
		s.addToken(TokMinus, nil)
	case '+':
		s.addToken(TokPlus, nil)
	case ';':
		s.addToken(TokSemicolon, nil)
	case '*':
		s.addToken(TokStar, nil)
	case '!':
		t := TokBang
		if s.match('=') {
			t = TokBangEqual
		}
		s.addToken(t, nil)
	case '=':
		t := TokEqual
		if s.match('=') {
			t = TokEqualEqual
		}
		s.addToken(t, nil)
	case '<':
		t := TokLess
		if s.match('=') {
			t = TokLessEqual
		}
		s.addToken(t, nil)
	case '>':
		t := TokGreater
		if s.match('=') {
			t = TokGreaterEqual
		}
		s.addToken(t, nil)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.Done() {
				s.advance()
			}
		} else {
			s.addToken(TokSlash, nil)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			handleErr(s.line, fmt.Sprintf("Unexpected character: %s", string(c)))
			return false
		}
	}
	return true
}

func (s *Scanner) advance() byte {
	c := s.src[s.current]
	s.current++
	return c
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	s.tokens = append(s.tokens, Token{
		Type:    tokenType,
		Lexeme:  s.src[s.start:s.current],
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) match(expected byte) bool {
	if s.Done() || s.src[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.Done() {
		return '\000'
	}
	return s.src[s.current]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.Done() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.Done() {
		handleErr(s.line, fmt.Sprintf("Incomplete string: %s", s.src[s.start:s.current]))
	}
	s.advance() // closing "
	value := s.src[s.start+1 : s.current-1]
	s.addToken(TokString, value)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	value, _ := strconv.ParseFloat(s.src[s.start:s.current], 64)
	s.addToken(TokNumber, value)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.src) {
		return '\000'
	}
	return s.src[s.current+1]
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	tokenType := TokIdentifier
	if keywordType, ok := keywords[s.src[s.start:s.current]]; ok {
		tokenType = keywordType
	}
	s.addToken(tokenType, nil)
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}
