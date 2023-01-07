package bittybox

// stack is a LIFO data structure
type stack struct {
	Values []token
}

// Pop removes the token at the top of the stack and returns its value
func (s *stack) Pop() token {
	end := len(s.Values) - 1
	if end < 0 {
		return token{}
	}
	token := s.Values[end]
	s.Values = s.Values[:end]
	return token
}

func (s *stack) Push(i token) {
	s.Values = append(s.Values, i)
}

// Peek returns the token at the top of the stack
func (s *stack) Peek() token {
	if len(s.Values) == 0 {
		return token{}
	}
	return s.Values[len(s.Values)-1]
}

// EmptyInto dumps all tokens from one stack to another
func (s *stack) EmptyInto(o *stack) {
	if !s.IsEmpty() {
		for i := s.Length() - 1; i >= 0; i-- {
			o.Push(s.Pop())
		}
	}
}

// IsEmpty checks if there are any tokens in the stack
func (s *stack) IsEmpty() bool {
	return len(s.Values) == 0
}

// Length returns the amount of tokens in the stack
func (s *stack) Length() int {
	return len(s.Values)
}
