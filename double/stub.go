package double

type Stub struct {
	RegisteredCalls []Call
}

func (s *Stub) On(methodName string, args ...interface{}) Call {
	call := Call{MethodName: methodName, Arguments: args}
	s.RegisteredCalls = append(s.RegisteredCalls, call)
	return call
}
