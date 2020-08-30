package parser

import "unicode/utf8"

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
	val int32
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

func NewNumberToken(p *Parser) *NumberToken {
	return &NumberToken{
		val: 0,
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
}

// NewParser get a Parser
func NewParser(input string) *Parser {
	return &Parser{
		input:       input,
		curTokenPos: 0,
	}
}

// Reset the input string
func (p *Parser) Reset(str string) {
	p.input = str
}

func (p *Parser) getNextToken() Token {
	if p.curTokenPos < len(p.input) {
		ch, w := utf8.DecodeRuneInString(p.input[p.curTokenPos:])

		if ch >= '0' && ch <= '9' {
			p.curTokenPos += w

		}

		if ch == '+' {
			return NewAddToken(p)
		}

		if ch == ' ' || ch == '\t' || ch == '\r' {

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
