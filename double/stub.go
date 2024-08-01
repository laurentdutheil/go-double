package double

import "github.com/stretchr/objx"

type Stub[T interface{}] struct {
	predefinedCalls Calls
	t               TestingT
	caller          *T
	testData        objx.Map
}

func (s *Stub[T]) On(methodName string, arguments ...interface{}) *Call {
	call := NewCall(methodName, arguments...)
	s.predefinedCalls = append(s.predefinedCalls, call)
	return call
}

func (s *Stub[T]) Called(arguments ...interface{}) Arguments {
	methodInformation := s.getMethodInformation()
	return s.MethodCalled(*methodInformation, arguments...)
}

func (s *Stub[T]) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	if s.t == nil {
		panic("Please use double.New constructor to initialize correctly.")
	}

	foundCall := s.predefinedCalls.find(methodInformation.Name, arguments...)

	if foundCall == noCallFound && methodInformation.NumOut > 0 {
		s.t.Errorf("I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", methodInformation.Name)
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

func (s *Stub[T]) TestData() objx.Map {
	if s.testData == nil {
		s.testData = make(objx.Map)
	}

	return s.testData
}

func (s *Stub[T]) When(method interface{}, arguments ...interface{}) *Call {
	functionName, err := GetFunctionName(method)
	if err != nil {
		panic("Please pass the function as an argument : stub.When(stub.Method)")
	}

	call := NewCall(functionName, arguments...)
	s.predefinedCalls = append(s.predefinedCalls, call)

	return call
}

func (s *Stub[T]) getMethodInformation() *MethodInformation {
	if s.t == nil {
		panic("Please use double.New constructor to initialize correctly.")
	}

	methodInformation, err := GetCallingMethodInformation(s.caller)
	if err != nil {
		s.t.Errorf(err.Error() + "\n\tUse MethodCalled instead of Called in stub implementation.")
		s.t.FailNow()
	}
	return methodInformation
}
