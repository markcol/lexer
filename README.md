[![Build Status](https://travis-ci.org/markcol/lexer.svg?branch=master)](https://travis-ci.org/markcol/lexer)
[![GoDoc](https://godoc.org/github.com/markcol/lexer?status.svg)](http://godoc.org/github.com/markcol/lexer)
# Lexer
A generic lexer as described by Rob Pike in his talk
"[Lexical Scanning in Go](http://cuddle.googlecode.com/hg/talk/lex.html)."
The code modified from the code presented in the talk in
order to make it usable as a separate package.
# Usage
```golang
import "github.com/markcol/lexer"

// Define token constants. The EOF token is -1, add your tokens after that.
const  (
    TokenEOF = lexer.TokenEOF + iota
    TokenSpace
    TokenIdent
	...
	)

// startState is the starting state of your lexer. The function should
// return a pointer to the next state function, or nil when parsing is
// complete.
func startState(l *lexer.Lexer) lexer.StateFn {
    ...
}

func main() {
    ...
    lex := NewLexer("Example", "Test\n", startState)
    for t := lex.NextToken() {
        switch t.Typ {
		case TokenEOF:
		  break;
        case TokenIdent:
        ...
		}
    }
}

```

