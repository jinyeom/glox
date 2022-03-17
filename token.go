package main

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

	// Keywords
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

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}
