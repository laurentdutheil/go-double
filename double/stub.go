package double

import "C"

type Stub[T interface{}] struct {
	predefinedCalls Calls
	t               TestingT
	caller          *T
}

func (s *Stub[T]) On(methodName string, arguments ...interface{}) *Call {
	call := NewCall(methodName, arguments...)
	s.predefinedCalls = append(s.predefinedCalls, call)
	return call
}

func (s *Stub[T]) Called(arguments ...interface{}) Arguments {
	if s.t == nil {
		panic("Please use double.New constructor to initialize correctly.")
	}

	method := GetCallingMethod(s.caller)
	return s.MethodCalled(method, arguments...)
}

func (s *Stub[T]) MethodCalled(method Method, arguments ...interface{}) Arguments {
	foundCall := s.predefinedCalls.find(method.Name, arguments...)

	if foundCall == noCallFound && method.NumOut > 0 {
		s.t.Errorf("I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", method.Name)
		s.t.FailNow()
	}

	return foundCall.called(arguments...)
}

func (s *Stub[T]) Test(t TestingT) {
	s.t = t
}

func (s *Stub[T]) PredefinedCalls() []*Call {
	return s.predefinedCalls
}
