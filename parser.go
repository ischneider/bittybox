package bittybox

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"text/scanner"
)

type parser struct {
	s    *scanner.Scanner
	err  error
	t    rune
	txt  string
	toks *stack
}

func newParser(in io.Reader) *parser {
	s := &scanner.Scanner{}
	s = s.Init(in)
	s.Mode = scanner.ScanFloats | scanner.ScanInts | scanner.ScanIdents
	p := &parser{s: s}
	s.Error = p.onError
	return p
}

func newParserString(expr string) *parser {
	return newParser(bytes.NewBufferString(expr))
}

func (p *parser) consume() error {
	r := p.s.Scan()
	if p.err != nil {
		return p.err
	}
	if r > 0 {
		valid := isOperator(r) || r == '(' || r == ')'
		if !valid {
			p.error("unexpected token: " + scanner.TokenString(r))
		}
	}
	p.t = r
	p.txt = p.s.TokenText()
	return p.err
}

func (p *parser) onError(s *scanner.Scanner, msg string) {
	pos := s.Position
	p.err = ErrInvalidSyntax{Message: msg, Line: pos.Line, Column: pos.Column}
}

func (p *parser) error(msg string, args ...interface{}) error {
	pos := p.s.Position
	return ErrInvalidSyntax{Message: fmt.Sprintf(msg, args...), Line: pos.Line, Column: pos.Column}
}

func (p *parser) Parse() (stack, error) {
	p.toks = &stack{}
	return *p.toks, p.parse()
}

func (p *parser) parse() error {
	err := p.consume()
	if err != nil {
		return err
	}
	err = p.expr()
	if err != nil {
		return err
	}
	return p.expect(scanner.EOF)
}

func (p *parser) expect(r rune) error {
	if p.t != r {
		return p.error("expected %s, got %s", scanner.TokenString(r), scanner.TokenString(p.t))
	}
	return nil
}

func (p *parser) tok(t tokenType, s string) {
	p.toks.Push(token{Type: t, Value: s})
}

func (p *parser) float() error {
	v, err := strconv.ParseFloat(p.txt, 64)
	if err != nil {
		return p.error("invalid float: %q", p.txt)
	}
	p.toks.Push(token{Type: float, Val: v})
	return p.consume()
}

func (p *parser) neg() error {
	p.tok(unary, "-")
	err := p.consume()
	if err != nil {
		return err
	}
	return p.unit()
}

func (p *parser) ident() error {
	ident := p.txt
	err := p.consume()
	if err != nil {
		return err
	}
	if p.t == '(' {
		if _, ok := funcs[ident]; !ok {
			return p.error("no function named: %q", ident)
		}
		p.tok(function, ident)
		return p.nested()
	}
	t := variable
	if _, ok := consts[ident]; ok {
		t = constant
	}
	p.tok(t, ident)
	return nil
}

func (p *parser) unit() error {
	var err error
	t := p.t
	switch p.t {
	case scanner.Float, scanner.Int:
		err = p.float()
	case '(':
		err = p.nested()
	case '-':
		err = p.neg()
	case scanner.Ident:
		err = p.ident()
	default:
		err = p.error("unexpected token: %s", scanner.TokenString(t))
	}
	return err
}

func (p *parser) nested() error {
	p.tok(lparen, "(")
	err := p.consume()
	if err != nil {
		return err
	}
	err = p.expr()
	if err != nil {
		return err
	}
	err = p.expect(')')
	if err != nil {
		return err
	}
	err = p.consume()
	if err != nil {
		return err
	}
	p.tok(rparen, ")")
	return nil
}

func (p *parser) expr() error {
	err := p.unit()
	for isOperator(p.t) && err == nil {
		p.tok(binary, p.txt)
		err = p.consume()
		if err == nil {
			err = p.unit()
		}
	}
	return err
}

func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '^'
}
