package double

type Stub struct {
	PredefinedCalls []*Call
	t               TestingT
	caller          interface{}
}

func (s *Stub) On(methodName string, arguments ...interface{}) *Call {
	call := NewCall(methodName, arguments...)
	s.PredefinedCalls = append(s.PredefinedCalls, call)
	return call
}

func (s *Stub) Called(arguments ...interface{}) Arguments {
	method := GetCallingMethod(s.caller)
	return s.MethodCalled(method, arguments...)
}

func (s *Stub) MethodCalled(method Method, arguments ...interface{}) Arguments {
	foundCall := s.findPredefinedCall(method.Name, arguments...)

	if foundCall == noCallFound && method.NumOut > 0 {
		s.t.Errorf("I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", method.Name)
		s.t.FailNow()
	}

	foundCall.incrementNumberOfCall()
	return foundCall.ReturnArguments
}

func (s *Stub) Test(t TestingT) {
	s.t = t
}

func (s *Stub) Caller(caller interface{}) {
	s.caller = caller
}

func (s *Stub) findPredefinedCall(methodName string, arguments ...interface{}) *Call {
	for _, predefinedCall := range s.PredefinedCalls {
		if methodName == predefinedCall.MethodName {
			if !predefinedCall.Arguments.Equal(arguments...) ||
				predefinedCall.alreadyCalledPredefinedTimes() {
				continue
			}
			return predefinedCall
		}
	}
	return noCallFound
}

var noCallFound = NewCall("-CallNotFound-")
