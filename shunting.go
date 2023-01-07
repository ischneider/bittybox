package bittybox

func shuntingYard(s stack, vars []string) (stack, error) {
	postfix := stack{}
	operators := stack{}
	idx := map[string]int{}
	for i, v := range vars {
		idx[v] = i
	}
	missing := map[string]bool{}
	for _, v := range s.Values {
		switch v.Type {
		case unary, function, lparen:
			operators.Push(v)
		case binary:
			for !operators.IsEmpty() {
				val := v.Value
				top := operators.Peek()
				if top.Type == unary || top.Type == function {
					postfix.Push(operators.Pop())
					continue
				}
				opr := oprData[val]
				topr := oprData[top.Value]
				if (opr.prec <= topr.prec && !opr.rAsoc) ||
					(opr.prec < topr.prec && opr.rAsoc) {
					postfix.Push(operators.Pop())
					continue
				}
				break
			}
			operators.Push(v)
		case rparen:
			for i := operators.Length() - 1; i >= 0; i-- {
				if operators.Values[i].Type != lparen {
					postfix.Push(operators.Pop())
					continue
				} else {
					operators.Pop()
					break
				}
			}
		case variable:
			key, ok := idx[v.Value]
			if !ok {
				missing[v.Value] = true
			}
			postfix.Push(token{variable, v.Value, float64(key), nil})
		default:
			postfix.Push(v)
		}
	}
	if len(missing) > 0 {
		names := make([]string, 0, len(missing))
		for i := range missing {
			names = append(names, i)
		}
		return postfix, ErrUnboundVars{names}
	}
	operators.EmptyInto(&postfix)
	return postfix, nil
}
