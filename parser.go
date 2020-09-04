package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

var eol = errors.New("end of line")

type Value interface{}

type Token interface {
	Led(leftValue Value) Value
	Nud() Value
	LeftBinding() int32
	RightBinding() int32
	Literal() string
}

type EOFToken struct {
}

func (t *EOFToken) Led(leftValue Value) Value {
	return leftValue
}

func (t *EOFToken) Nud() Value {
	return 0
}

func (t *EOFToken) LeftBinding() int32 {
	return 0
}

func (t *EOFToken) RightBinding() int32 {
	return 0
}

func (t *EOFToken) Literal() string {
	return "eof"
}

var _ = Token(&EOFToken{})

var eof = &EOFToken{}

type NumberToken struct {
	val int64
	p   *Parser
}

func (n *NumberToken) Led(leftValue Value) Value {
	return 0
}

// Nud return a node with new Value
func (n *NumberToken) Nud() Value {
	return n.val
}

func (n *NumberToken) LeftBinding() int32 {
	return 100
}
func (n *NumberToken) RightBinding() int32 {
	return 100
}

func (n *NumberToken) Literal() string {
	return strconv.FormatInt(n.val, 10)
}

func NewNumberToken(p *Parser, val int64) *NumberToken {
	return &NumberToken{
		val: val,
		p:   p,
	}
}

// Number is Token
var _ = Token(&NumberToken{})

type AddToken struct {
	left  Value
	right Value
	lbp   int32
	rbp   int32
	p     *Parser
}

func NewAddToken(p *Parser) *AddToken {
	return &AddToken{
		lbp: 100,
		rbp: 10,
		p:   p,
	}
}

func (t *AddToken) LeftBinding() int32 {
	return t.lbp
}

func (t *AddToken) RightBinding() int32 {
	return t.rbp
}

func (t *AddToken) Led(leftValue Value) Value {
	t.left = leftValue
	t.right = t.p.Parse(t.RightBinding())

	var left int64
	var right int64

	left, _ = t.left.(int64)
	right, _ = t.right.(int64)
	fmt.Println("...left...+ right..", left+right)
	return left + right
}

func (t *AddToken) Nud() Value {
	t.left = t.p.Parse(t.LeftBinding())
	t.right = 0

	if val, ok := t.left.(int64); ok {
		return val
	}
	return 0
}

func (t *AddToken) Literal() string {
	return "+"
}

// Parser to parse expression
type Parser struct {
	input       string
	curTokenPos int
	previousPos int
}

// NewParser get a Parser

func NewParser(input string) *Parser {
	return &Parser{
		input:       input,
		curTokenPos: 0,
		previousPos: -1,
	}
}

// Reset the input string
func (p *Parser) Reset(str string) {
	p.input = str
}

func (p *Parser) consumeChar() (rune, error) {
	var r rune
	if p.curTokenPos < len(p.input) {
		r, w := utf8.DecodeRuneInString(p.input[p.curTokenPos:])
		p.previousPos = p.curTokenPos
		p.curTokenPos += w
		return r, nil
	}
	// 表示da
	return r, eol

}

func (p *Parser) preChar() (rune, error) {
	var r rune
	if p.previousPos < len(p.input) {
		r, _ := utf8.DecodeLastRuneInString(p.input[p.previousPos:])
		return r, nil
	}
	return r, eol
}

func (p *Parser) peekChar() (rune, error) {
	var r rune
	if p.curTokenPos < len(p.input) {
		r, _ := utf8.DecodeRuneInString(p.input[p.curTokenPos:])
		return r, nil
	}
	return r, eol
}

func (p *Parser) skipWhiteSpace() {
	for {
		ch, err := p.peekChar()

		if err == eol {
			return
		}
		isBlank := false
		switch ch {
		case '\n':
		case '\r':
		case '\t':
		case ' ':
			p.consumeChar()
			isBlank = true
		}
		if !isBlank {
			break
		}
	}
}

func (p *Parser) isDigit(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

func (p *Parser) getNextToken() Token {
	p.skipWhiteSpace()
	r, err := p.peekChar()
	if err == eol {
		return eof
	}
	switch r {
	case '+':
		p.consumeChar()
		return NewAddToken(p)
	default:
		var buffer bytes.Buffer
		if p.isDigit(r) {
			for p.isDigit(r) && err == nil {
				buffer.WriteRune(r)
				p.consumeChar()
				r, err = p.peekChar()
			}
			num := buffer.String()
			n, _ := strconv.ParseInt(num, 10, 32)
			return NewNumberToken(p, n)
		} else {

		}

	}
	return eof
}

// Parse parse expresion
func (p *Parser) Parse(rbp int32) Value {
	token := p.getNextToken()
	fmt.Printf("first token %s\n", token.Literal())
	old := token
	token = p.getNextToken()
	fmt.Printf("second token %s\n", token.Literal())
	leftValue := old.Nud()
	fmt.Printf("rbp  %d\n", rbp)

	for rbp < token.LeftBinding() && token != eof {
		old = token
		leftValue = old.Led(leftValue)
		token = p.getNextToken()
	}
	return leftValue
}
