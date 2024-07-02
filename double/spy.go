package double

type Spy struct {
	Stub
	ActualCalls []Call
}

func (s *Spy) Called(arguments ...interface{}) Arguments {
	functionName := getCallingFunctionName()
	return s.MethodCalled(functionName, arguments...)
}

func (s *Spy) MethodCalled(methodName string, arguments ...interface{}) Arguments {
	call := *NewCall(methodName, arguments...)
	s.ActualCalls = append(s.ActualCalls, call)

	return s.Stub.MethodCalled(methodName, arguments...)
}
