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

func (s *Spy) NumberOfCall(methodName string) int {
	count := 0
	for _, call := range s.ActualCalls {
		if call.MethodName == methodName {
			count++
		}
	}
	return count
}
