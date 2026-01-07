package bittybox_test

import (
	"fmt"

	"github.com/ischneider/bittybox"
)

func ExampleEvaluate() {
	fmt.Println(bittybox.Evaluate("5 + y", bittybox.Var{"y", 1}))
	// Output: 6 <nil>
}

func ExampleCompileExpr() {
	expr, err := bittybox.CompileExpr("5 + y", "y", "z")
	if err != nil {
		panic(err)
	}
	fmt.Println(expr.Evaluate([]float64{10}), expr.Vars)
	// Output: 15 [y]
}
