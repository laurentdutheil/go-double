package double

import (
	"github.com/stretchr/objx"
)

// Stub provides prepared answers to calls made during test.
// For an example of its usage, refer to the "Example Usage" section at the top
// of this document.
type Stub struct {
	predefinedCalls Calls
	t               TestingT
	caller          interface{}
	testData        objx.Map
}

// On starts a description of an expectation of the specified method
// being called.
//
//	Stub.On("Method", arg1, arg2)
func (s *Stub) On(methodName string, arguments ...interface{}) *Call {
	call := NewCall(methodName, arguments...)
	s.predefinedCalls = append(s.predefinedCalls, call)
	return call
}

// Called tells the stub object that a method has been called, and gets an array
// of arguments to return.  Fail the test if the call is unexpected (i.e. not preceded by
// appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (s *Stub) Called(arguments ...interface{}) Arguments {
	methodInformation := s.getMethodInformation()
	return s.MethodCalled(*methodInformation, arguments...)
}

// MethodCalled tells the stub object that a method has been called, and gets an array
// of arguments to return.  Fail the test if the call is unexpected (i.e. not preceded by
// appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (s *Stub) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	s.checkInitialization()

	foundCall := s.predefinedCalls.find(methodInformation.Name, arguments...)

	if foundCall == noCallFound && methodInformation.NumOut > 0 {
		s.t.Errorf("I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", methodInformation.Name)
		s.t.FailNow()
	}

	return foundCall.called(arguments...)
}

// Test sets the test struct variable of the stub object.
// If you don't use the double.New constructor, you have to set it yourself.
func (s *Stub) Test(t TestingT) {
	s.t = t
}

// Caller sets the caller struct (the stub object itself).
// If you don't use the double.New constructor, you have to set it yourself.
func (s *Stub) Caller(caller interface{}) {
	s.caller = caller
}

// PredefinedCalls return the predefined calls of the Stub
func (s *Stub) PredefinedCalls() []*Call {
	return s.predefinedCalls
}

// TestData holds any data that might be useful for testing.  Testify ignores
// this data completely allowing you to do whatever you like with it.
func (s *Stub) TestData() objx.Map {
	if s.testData == nil {
		s.testData = make(objx.Map)
	}

	return s.testData
}

// When is similar of the On method, except you pass the method instead of the name.
// Panics if the method argument is not reflect.Func type.
//
//	Stub.When(Stub.Method, arg1, arg2)
func (s *Stub) When(method interface{}, arguments ...interface{}) *Call {
	s.checkInitialization()

	functionName, err := GetFunctionName(method)
	if err != nil {
		s.t.Errorf("Please pass the function as an argument : stub.When(stub.Method)")
		s.t.FailNow()
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
