package bittybox

import (
	"math"
	"strconv"
)

var oprData = map[string]struct {
	prec  int
	rAsoc bool // true = right // false = left
}{
	"^": {4, true},
	"*": {3, false},
	"/": {3, false},
	"+": {2, false},
	"-": {2, false},
}

var funcs = map[string]func(x float64) float64{
	"ln":    math.Log,
	"abs":   math.Abs,
	"cos":   math.Cos,
	"sin":   math.Sin,
	"tan":   math.Tan,
	"acos":  math.Acos,
	"asin":  math.Asin,
	"atan":  math.Atan,
	"sqrt":  math.Sqrt,
	"cbrt":  math.Cbrt,
	"ceil":  math.Ceil,
	"floor": math.Floor,
}

var consts = map[string]float64{
	"E":       math.E,
	"PI":      math.Pi,
	"PHI":     math.Phi,
	"SQRT2":   math.Sqrt2,
	"SQRTE":   math.SqrtE,
	"SQRTPI":  math.SqrtPi,
	"SQRTPHI": math.SqrtPhi,
}

func solve(tokens stack, fstack []float64, vars []float64) (float64, []float64) {
	fstack = fstack[:0]
	for _, v := range tokens.Values {
		switch v.Type {
		case variable:
			fstack = append(fstack, vars[int(v.Val)])
		case float:
			fstack = append(fstack, v.FVal())
		case function:
			top := len(fstack) - 1
			fstack[top] = funcs[v.Value](fstack[top])
		case constant:
			if val, ok := consts[v.Value]; ok {
				fstack = append(fstack, val)
			} else {
				panic(v.Value)
			}
		case binary:
			top := len(fstack) - 1
			next := top - 1
			y, x := fstack[top], fstack[next]
			var result float64
			switch v.Value {
			case "^":
				result = math.Pow(x, y)
			case "+":
				result = x + y
			case "-":
				result = x - y
			case "*":
				result = x * y
			case "/":
				result = x / y
			}
			fstack[next] = result
			fstack = fstack[:top]
		case unary:
			top := len(fstack) - 1
			fstack[top] = -fstack[top]
		default:
			panic("not expected/handled " + strconv.Itoa(int(v.Type)))
		}
	}
	return fstack[0], fstack
}
