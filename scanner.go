package main

import (
	"fmt"
	"strconv"
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

type Scanner struct {
	src                 string
	start, offset, line int
	tokens              []Token
}

func NewScanner(src string) *Scanner {
	return &Scanner{
		src:    src,
		start:  0,
		offset: 0,
		line:   1,
		tokens: make([]Token, 0),
	}
}

func (s *Scanner) Scan() ([]Token, error) {
	for !s.done() {
		if err := s.scanToken(); err != nil {
			return nil, err
		}
		s.start = s.offset
	}
	s.tokens = append(s.tokens, Token{TokEOF, "", nil, s.line})
	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	c := s.step()
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
			for s.peek() != '\n' && !s.done() {
				s.step()
			}
		} else {
			s.addToken(TokSlash, nil)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		return s.string()
	default:
		switch {
		case isDigit(c):
			s.number()
		case isAlpha(c):
			s.identifier()
		default:
			return fmt.Errorf("unexpected character: %s", string(c))
		}
	}
	return nil
}

func (s *Scanner) step() byte {
	c := s.src[s.offset]
	s.offset++
	return c
}

func (s *Scanner) done() bool {
	return s.offset >= len(s.src)
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	s.tokens = append(s.tokens, Token{
		Type:    tokenType,
		Lexeme:  s.src[s.start:s.offset],
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) match(expected byte) bool {
	if s.done() || s.src[s.offset] != expected {
		return false
	}
	s.offset++
	return true
}

func (s *Scanner) peek() byte {
	if s.done() {
		return '\000'
	}
	return s.src[s.offset]
}

func (s *Scanner) peekNext() byte {
	if s.offset+1 >= len(s.src) {
		return '\000'
	}
	return s.src[s.offset+1]
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.done() {
		if s.peek() == '\n' {
			s.line++
		}
		s.step()
	}
	if s.done() {
		return fmt.Errorf("incomplete string: %s", s.src[s.start:s.offset])
	}
	s.step() // closing "
	value := s.src[s.start+1 : s.offset-1]
	s.addToken(TokString, value)
	return nil
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.step()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.step()
		for isDigit(s.peek()) {
			s.step()
		}
	}
	value, _ := strconv.ParseFloat(s.src[s.start:s.offset], 64)
	s.addToken(TokNumber, value)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.step()
	}
	tokenType := TokIdentifier
	if keywordType, ok := keywords[s.src[s.start:s.offset]]; ok {
		tokenType = keywordType
	}
	s.addToken(tokenType, nil)
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
