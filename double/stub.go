package double

import (
	"fmt"
)

type Stub struct {
	PredefinedCalls []*Call
}

func (s *Stub) On(methodName string, arguments ...interface{}) *Call {
	call := NewCall(methodName, arguments...)
	s.PredefinedCalls = append(s.PredefinedCalls, call)
	return call
}

func (s *Stub) Called(caller interface{}, arguments ...interface{}) Arguments {
	method := GetCallingMethod(caller)
	return s.MethodCalled(method, arguments...)
}

func (s *Stub) MethodCalled(method Method, arguments ...interface{}) Arguments {
	numberOfReturnArguments := method.NumOut
	if numberOfReturnArguments == 0 {
		return nil
	}

	foundCall := s.findPredefinedCall(method.Name, arguments...)
	if foundCall == nil {
		errorMessage := fmt.Sprintf("I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", method.Name)
		panic(errorMessage)
	}

	return foundCall.ReturnArguments
}

func (s *Stub) findPredefinedCall(methodName string, arguments ...interface{}) *Call {

	for _, predefinedCall := range s.PredefinedCalls {
		if methodName == predefinedCall.MethodName {
			if !predefinedCall.Arguments.Diff(arguments...) {
				return nil
			}
			return predefinedCall
		}
	}
	return nil
}
