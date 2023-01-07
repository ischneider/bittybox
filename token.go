package bittybox

import (
	"fmt"
)

type tokenType int

type token struct {
	Type  tokenType
	Value string
	Val   float64
	F     func() float64
}

func (t token) String() string {
	switch t.Type {
	case constant, variable, function, binary, lparen, rparen:
		return t.Value
	case unary:
		return "neg"
	case float:
		return fmt.Sprint(t.Val)
	}
	return fmt.Sprintf("<%s %s %f>", tokenTypeName(t.Type), t.Value, t.Val)
}

func (t token) FVal() float64 {
	if t.F != nil {
		return t.F()
	}
	return t.Val
}

const (
	constant tokenType = iota
	lparen
	rparen
	function
	binary
	unary
	float
	variable
)

func tokenTypeName(t tokenType) string {
	switch t {
	case lparen:
		return "("
	case rparen:
		return ")"
	case constant:
		return "const"
	case function:
		return "func"
	case binary:
		return "binary"
	case unary:
		return "unary"
	case float:
		return "float"
	case variable:
		return "var"
	}
	panic(t)
}
