package bittybox

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"
)

var basicTests = []struct {
	expr   string
	expect float64
}{
	{"5", 5},
	{"5.5 + 5.5", 5.5 + 5.5},
	{"5 + 10", 5 + 10},
	{"5 + 10^2", 5 + 10*10},
	{"5 + 10 * 5", 5 + 10*5},
	{"5 + 10 / 2", 5 + 10/2},
	{"-5 + 10 / 2 * 9 - 5 ^ -2", -5 + 10/2*9 - math.Pow(5, -2)},
	{"(5 + 5) / 2", (5 + 5) / 2},
	{"(5 + (5 * 5))", (5 + (5 * 5))},
	{"(5 + 5) * (5 / 1) + ((5 / 1) * 5 * (5+5))", (5+5)*(5/1) + ((5 / 1) * 5 * (5 + 5))},
	{"(5 + (5 * 5))", (5 + (5 * 5))},
	{"sqrt(PI)", math.Sqrt(math.Pi)},
	{"sqrt(10)", math.Sqrt(10)},
	{"sqrt(sin(0) + cos(1))", math.Sqrt(math.Sin(0) + math.Cos(1))},
	{"sqrt(5 * 10)", math.Sqrt(5 * 10)},
	{"-5 + 10", -5 + 10},
	{"-5 + -5", -5 + -5},
	{"-5 - -5", 0},
	{"1 / 0", math.Inf(1)},
	{"-sqrt(5)", -math.Sqrt(5)},
	{"cos(-1)", math.Cos(-1)},
}

var badSyntaxTests = []struct {
	expr string
	err  string
}{
	{"", `0:0 : unexpected token: EOF`},
	{"sqrt(", `1:6 : unexpected token: EOF`},
	{"sqrt&", `1:5 : expected EOF, got "&"`},
	{"(0", `1:3 : expected ")", got EOF`},
	{"0xf", `1:1 : invalid float: "0xf"`},
	{"0e", "1:1 : exponent has no digits"},
	{"0N", "1:2 : expected EOF, got Ident"},
	{"5 *", `1:4 : unexpected token: EOF`},
	{"func(5)", `1:5 : no function named: "func"`},
	{"-I\xfa", `1:2 : invalid UTF-8 encoding`},
	{"a \xfa", `0:1 : invalid UTF-8 encoding`},
	{"()\xfa", `1:2 : invalid UTF-8 encoding`},
	{"sqrt()", `1:6 : unexpected token: ")"`},
}

func mustEvaluate(e string, vars []Var) float64 {
	f, err := Evaluate(e, vars...)
	if err != nil {
		panic(err)
	}
	return f
}

func FuzzEm(f *testing.F) {
	for _, t := range basicTests {
		f.Add(t.expr)
	}
	for _, t := range badSyntaxTests {
		f.Add(t.expr)
	}
	f.Fuzz(func(t *testing.T, e string) {
		s, err := newParserString(e).Parse()
		if err == nil {
			s, err = shuntingYard(s, nil)
			if err == nil {
				solve(s, nil, nil)
			}
		}
	})
}

func TestBadSyntax(t *testing.T) {
	for _, c := range badSyntaxTests {
		s, err := newParserString(c.expr).Parse()
		if err == nil {
			t.Errorf("%q : want err, got nil - stack is %v", c.expr, s)
		}
		if got := err.Error(); got != c.err {
			t.Errorf("[%q] want %s, got %s - stack is %v", c.expr, c.err, got, s)
		}
	}
}

func TestMissingVars(t *testing.T) {
	s, err := newParserString("x + 5").Parse()
	if err != nil {
		panic(err)
	}
	_, err = shuntingYard(s, nil)
	if err == nil {
		t.Error("expected error")
	}
	if err.Error() != "missing variables: [x]" {
		t.Errorf("want x, got %s", err)
	}
}

func TestNonVars(t *testing.T) {
	for _, c := range basicTests {
		answer := mustEvaluate(c.expr, nil)
		if answer != c.expect {
			t.Errorf("expected %f, got %f", c.expect, answer)
		}
	}
}

func TestVars(t *testing.T) {
	if f := mustEvaluate("x + y", []Var{{Name: "x", Value: 1}, {Name: "y", Value: 2}}); f != 3 {
		t.Error("expected 3")
	}
}

func TestIt(t *testing.T) {
	s := "(x + y) / (x - y)"
	p := newParser(strings.NewReader(s))
	stack, err := p.Parse()
	if err != nil {
		panic(err)
	}
	stack, err = shuntingYard(stack, []string{"x", "y"})
	if err != nil {
		panic(err)
	}
	st := time.Now()
	var fstack []float64
	vars := []float64{5, 10}
	for i, ii := 0, 256*256; i < ii; i++ {
		var answer float64
		answer, fstack = solve(stack, fstack, vars)
		if answer != -3 {
			t.Fail()
		}
	}
	fmt.Println("process loop", time.Since(st))
}

func BenchmarkIt(b *testing.B) {
	s := "x + y + z"
	p := newParser(strings.NewReader(s))
	stack, err := p.Parse()
	if err != nil {
		panic(err)
	}
	buf := make([]float64, len(stack.Values))
	stack, err = shuntingYard(stack, []string{"x", "y", "z"})
	if err != nil {
		panic(err)
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		fi := float64(i)
		vars := []float64{fi, fi, fi}
		var answer float64
		answer, buf = solve(stack, buf, vars)
		if answer != fi*3 {
			fmt.Println(answer)
			b.Fail()
		}
	}
}
