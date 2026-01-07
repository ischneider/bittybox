package bittybox

import "fmt"

// ErrInvalidSyntax describes an invalid expression.
type ErrInvalidSyntax struct {
	Message string
	Line    int
	Column  int
}

func (e ErrInvalidSyntax) Error() string {
	return fmt.Sprintf("%d:%d : %s", e.Line, e.Column, e.Message)
}

// ErrUnboundVars is returned when an expression is missing variables.
type ErrUnboundVars struct {
	Missing []string
}

func (e ErrUnboundVars) Error() string {
	return fmt.Sprintf("missing variables: %v", e.Missing)
}

type Float interface {
	float32 | float64
}

// Var is a name, value pair.
type Var struct {
	Name  string
	Value float64
}

// Expr is a compiled expression that can be reused serially.
type Expr struct {
	s stack
	f []float64
	// Vars are the variables referenced in this Expr
	Vars []string
}

// Compute evaluates the expression with the provided values.
// Values must be provided in the order of their corresponding varNames.
func (e *Expr) Evaluate(vals []float64) float64 {
	r, f := solve(e.s, e.f, vals)
	e.f = f
	return r
}

// CompileExpr compiles the provided forumula into an Expr.
// VarNames are the allowed set of variables for evaluation.
func CompileExpr(formula string, varNames ...string) (Expr, error) {
	s, err := newParserString(formula).Parse()
	if err != nil {
		return Expr{}, err
	}
	s, err = shuntingYard(s, varNames)
	if err != nil {
		return Expr{}, err
	}
	var vars []string
	for _, t := range s.Values {
		if t.Type == variable {
			vars = append(vars, t.Value)
		}
	}
	return Expr{s, nil, vars}, err
}

// Compute will compile and evaluate an expression.
func Evaluate(formula string, vars ...Var) (float64, error) {
	names := make([]string, len(vars))
	vals := make([]float64, len(vars))
	for i := range vars {
		names[i] = vars[i].Name
		vals[i] = vars[i].Value
	}
	expr, err := CompileExpr(formula, names...)
	if err != nil {
		return 0, err
	}
	return expr.Evaluate(vals), nil
}
