package double

type Spy struct {
	Stub
	ActualCalls []Call
}

func (s *Spy) Called(caller interface{}, arguments ...interface{}) Arguments {
	method := GetCallingMethod(caller)
	return s.MethodCalled(method, arguments...)
}

func (s *Spy) MethodCalled(method Method, arguments ...interface{}) Arguments {
	call := *NewCall(method.Name, arguments...)
	s.ActualCalls = append(s.ActualCalls, call)

	return s.Stub.MethodCalled(method, arguments...)
}

func (s *Spy) NumberOfCalls(methodName string) int {
	count := 0
	for _, call := range s.ActualCalls {
		if call.MethodName == methodName {
			count++
		}
	}
	return count
}

func (s *Spy) NumberOfCallsWithArguments(methodName string, arguments ...interface{}) int {
	count := 0
	for _, call := range s.ActualCalls {
		if call.MethodName == methodName && call.Arguments.Equal(arguments...) {
			count++
		}
	}
	return count
}
