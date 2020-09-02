package parser

import (
	"bytes"
	"errors"
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

var _ = Token(&EOFToken{})

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

	var left int32
	var right int32

	left, _ = t.left.(int32)
	right, _ = t.right.(int32)
	return left + right
}

func (t *AddToken) Nud() Value {
	t.left = t.p.Parse(t.LeftBinding())
	t.right = 0

	if val, ok := t.left.(int32); ok {
		return val
	}
	return 0
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
	if p.previousPos < len(p.input) {
		r, _ := utf8.DecodeLastRuneInString(p.input[p.curTokenPos:])
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
		switch ch {
		case '\n':
		case '\r':
		case '\t':
		case ' ':
			p.consumeChar()
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
		return &EOFToken{}
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
	return &EOFToken{}
}

// Parse parse expresion
func (p *Parser) Parse(rbp int32) Value {
	token := p.getNextToken()
	old := token
	token = p.getNextToken()
	leftValue := old.Nud()

	for rbp < token.LeftBinding() {
		old = token
		token = p.getNextToken()
		leftValue = old.Led(leftValue)
	}
	return leftValue
}
