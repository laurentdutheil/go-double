package double

type Spy[T interface{}] struct {
	Stub[T]
	ActualCalls []Call
}

func (s *Spy[T]) Called(arguments ...interface{}) Arguments {
	if s.t == nil {
		panic("Please use double.New constructor to initialize correctly.")
	}

	method := GetCallingMethod(s.caller)
	return s.MethodCalled(method, arguments...)
}

func (s *Spy[T]) MethodCalled(method Method, arguments ...interface{}) Arguments {
	call := *NewCall(method.Name, arguments...)
	s.ActualCalls = append(s.ActualCalls, call)

	return s.Stub.MethodCalled(method, arguments...)
}

func (s *Spy[T]) NumberOfCalls(methodName string) int {
	count := 0
	for _, call := range s.ActualCalls {
		if call.MethodName == methodName {
			count++
		}
	}
	return count
}

func (s *Spy[T]) NumberOfCallsWithArguments(methodName string, arguments ...interface{}) int {
	count := 0
	for _, call := range s.ActualCalls {
		if call.MethodName == methodName && call.Arguments.Equal(arguments...) {
			count++
		}
	}
	return count
}
