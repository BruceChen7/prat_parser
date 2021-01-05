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

type NudFunc func(t Token) Value
type LedFunc func(t Token, leftValue Value) Value

type BinaryToken struct {
	rightBinding int32
	leftBinding  int32
	token        string
	nudFunc      NudFunc
	ledFunc      LedFunc
	p            *Parser
}

func (b BinaryToken) LeftBinding() int32 {
	return b.leftBinding
}

func (b BinaryToken) RightBinding() int32 {
	return b.rightBinding
}

func (b BinaryToken) Literal() string {
	return b.token
}

func (b BinaryToken) Nud() Value {
	return b.nudFunc(b)
}

func (b BinaryToken) Led(leftValue Value) Value {
	return b.ledFunc(b, leftValue)
}

//
var _ = Token(&BinaryToken{})

func NewBinaryToken(l, r int32, token string, led LedFunc, nud NudFunc, p *Parser) *BinaryToken {
	return &BinaryToken{
		rightBinding: r,
		leftBinding:  l,
		token:        token,
		nudFunc:      nud,
		ledFunc:      led,
		p:            p,
	}
}

func NewAddToken(p *Parser) *BinaryToken {
	return NewBinaryToken(100, 10, "+", AddTokenLed, AddTokenNud, p)
}

func AddTokenLed(t Token, leftValue Value) Value {
	if token, ok := t.(BinaryToken); ok {
		rightValue := token.p.Parse(token.RightBinding())

		left, _ := leftValue.(int64)
		right, _ := rightValue.(int64)
		fmt.Println("...left...+ right..", left+right)
		return left + right
	}
	return 0
}

func AddTokenNud(t Token) Value {
	if token, ok := t.(BinaryToken); ok {
		leftVal := token.p.Parse(t.LeftBinding())

		if val, ok := leftVal.(int64); ok {
			return val
		}
	}
	return 0
}

func MinusTokenLed(token Token, leftVal Value) Value {
	if m, ok := token.(BinaryToken); ok {
		rightVal := m.p.Parse(m.RightBinding())
		var left int64
		var right int64
		left, _ = leftVal.(int64)
		right, _ = rightVal.(int64)

		return left - right
	}
	return 0
}

func MinusTokenNud(token Token) Value {
	if m, ok := token.(BinaryToken); ok {
		right := m.p.Parse(m.LeftBinding())
		r, _ := right.(int64)
		return -1 * r
	}
	return 0
}

func NewMinsToken(p *Parser) *BinaryToken {
	return &BinaryToken{
		rightBinding: 10,
		leftBinding:  100,
		token:        "-",
		nudFunc:      MinusTokenNud,
		ledFunc:      MinusTokenLed,
		p:            p,
	}
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
	case '-':
		p.consumeChar()
		return NewMinsToken(p)
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

var token Token

// Parse parse expresion
func (p *Parser) Parse(rbp int32) Value {
	old := token
	token = p.getNextToken()
	fmt.Printf("second token %s\n", token.Literal())
	leftValue := old.Nud()
	fmt.Printf("rbp  %d\n", rbp)
	for rbp < token.LeftBinding() && token != eof {
		old = token
		// 这里必须先使用token
		token = p.getNextToken()
		leftValue = old.Led(leftValue)
	}
	return leftValue
}

func (p *Parser) Expr() Value {
	token = p.getNextToken()
	return p.Parse(0)
}
