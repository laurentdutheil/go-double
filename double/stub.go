package double

import (
	"github.com/stretchr/objx"
)

type Stub struct {
	predefinedCalls Calls
	t               TestingT
	caller          interface{}
	testData        objx.Map
}

func (s *Stub) On(methodName string, arguments ...interface{}) *Call {
	call := NewCall(methodName, arguments...)
	s.predefinedCalls = append(s.predefinedCalls, call)
	return call
}

func (s *Stub) Called(arguments ...interface{}) Arguments {
	methodInformation := s.getMethodInformation()
	return s.MethodCalled(*methodInformation, arguments...)
}

func (s *Stub) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	s.checkInitialization()

	foundCall := s.predefinedCalls.find(methodInformation.Name, arguments...)

	if foundCall == noCallFound && methodInformation.NumOut > 0 {
		s.t.Errorf("I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", methodInformation.Name)
		s.t.FailNow()
	}

	return foundCall.called(arguments...)
}

func (s *Stub) Test(t TestingT) {
	s.t = t
}

func (s *Stub) Caller(caller interface{}) {
	s.caller = caller
}

func (s *Stub) PredefinedCalls() []*Call {
	return s.predefinedCalls
}

func (s *Stub) TestData() objx.Map {
	if s.testData == nil {
		s.testData = make(objx.Map)
	}

	return s.testData
}

func (s *Stub) When(method interface{}, arguments ...interface{}) *Call {
	functionName, err := GetFunctionName(method)
	if err != nil {
		panic("Please pass the function as an argument : stub.When(stub.Method)")
	}

	call := NewCall(functionName, arguments...)
	s.predefinedCalls = append(s.predefinedCalls, call)

	return call
}

func (s *Stub) checkInitialization() {
	if s.t == nil || s.caller == nil {
		panic("Please use double.New constructor to initialize correctly.")
	}
}

func (s *Stub) getMethodInformation() *MethodInformation {
	s.checkInitialization()

	methodInformation, err := GetCallingMethodInformation(s.caller)
	if err != nil {
		s.t.Errorf(err.Error() + "\n\tUse MethodCalled instead of Called in stub implementation.")
		s.t.FailNow()
	}
	return methodInformation
}

type IStub interface {
	On(methodName string, arguments ...interface{}) *Call
	Called(arguments ...interface{}) Arguments
	MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments
	PredefinedCalls() []*Call
	TestData() objx.Map
	When(method interface{}, arguments ...interface{}) *Call
}

// Check if Stub implements all methods of IStub
var _ IStub = (*Stub)(nil)
