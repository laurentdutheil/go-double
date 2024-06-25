package double

type Stub struct {
	RegisteredCalls []Call
}

func (s *Stub) On(methodName string) Call {
	call := Call{MethodName: methodName}
	s.RegisteredCalls = append(s.RegisteredCalls, call)
	return call
}
