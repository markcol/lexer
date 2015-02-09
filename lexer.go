// Package lexer implements a generic lexer as described by Rob Pike
// in his talk "Lexical Scanning in Go." The package has been modified
// from the code presented in the talk in order to make it usable as a
// separate package.
package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// TokenType identifies the type used to represent lexical tokens.
type TokenType int

const (
	TokenError TokenType = -2 // error occured, value is text of error
	TokenEOF   TokenType = -1 // end of file token
)

// Token represents a token returned from the lexical scanner.
type Token struct {
	Typ TokenType // Type, such as itemNumber
	Val string    // Value, such as "23.2"
	Pos int       // location of token in input
}

const EOF = -1 // Rune returned to indicate EOF

func (i Token) String() string {
	switch i.Typ {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return i.Val
	}
	if len(i.Val) > 10 {
		return fmt.Sprintf("%.10q...", i.Val)
	}
	return fmt.Sprintf("%q", i.Val)
}

// StateFn represents the state of the scanner as a function that
// returns the next state.
type StateFn func(*Lexer) StateFn

// lexer holds the state of the scanner.
type Lexer struct {
	name    string     // used only for error reports
	input   string     // the string being scanned
	state   StateFn    // the next lexing function to enter
	Start   int        // start position of this item
	Pos     int        // current position in the input
	lastPos int        // position of last token in input
	Width   int        // width of last run from input
	tokens  chan Token // channel of scanned tokens
}

// NewLexer creates a new scanner for the input string.
func NewLexer(name, input string, startState StateFn) *Lexer {
	l := &Lexer{
		name:   name,
		input:  input,
		state:  startState,
		tokens: make(chan Token, 2), // two items sufficient
	}
	go l.run()
	return l
}

// Run lexes the input by execute state functions until the state is nil.
func (l *Lexer) run() {
	for state := l.state; state != nil; {
		state = state(l)
	}
}

// LineNumber returns the line number of the current position within the input string.
func (l *Lexer) LineNumber() int {
	return strings.Count(l.input[:l.lastPos], "\n") + 1
}

// NextToken returns the next item from the input.
func (l *Lexer) NextToken() Token {
	for {
		select {
		case token := <-l.tokens:
			l.lastPos = token.Pos
			return token
		default:
			l.state = l.state(l)
		}
	}
	panic("not reached")
}

// Emit passes an item back to the client
func (l *Lexer) Emit(t TokenType) {
	l.tokens <- Token{t, l.input[l.Start:l.Pos], l.Start}
	l.Start = l.Pos
}

// Next returns the next rune in the input.
func (l *Lexer) Next() rune {
	if l.Pos >= len(l.input) {
		l.Width = 0
		return EOF
	}
	r, w := utf8.DecodeRuneInString(l.input[l.Pos:])
	l.Width = w
	l.Pos += l.Width
	return r
}

// Ignore skips over the pending input before this point.
func (l *Lexer) Ignore() {
	l.Start = l.Pos
}

// Backup steps back one rune.
// Can be called only once per call of next.
func (l *Lexer) Backup() {
	l.Pos -= l.Width
}

// Peek returns but does not consume
// the next rune in the input.
func (l *Lexer) Peek() rune {
	r := l.Next()
	l.Backup()
	return r
}

// Accept consumes the next rune
// if it's from the valid set.
func (l *Lexer) Accept(valid string) bool {
	if strings.IndexRune(valid, l.Next()) >= 0 {
		return true
	}
	l.Backup()
	return false
}

// AcceptRun consumes a run of runes from the valid set.
func (l *Lexer) AcceptRun(valid string) {
	for strings.IndexRune(valid, l.Next()) >= 0 {
	}
	l.Backup()
}

// Errorf returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *Lexer) Errorf(format string, args ...interface{}) StateFn {
	l.tokens <- Token{
		TokenError,
		fmt.Sprintf(format, args...),
		l.Start,
	}
	return nil
}
