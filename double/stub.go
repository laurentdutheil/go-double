package double

import (
	"fmt"
	"reflect"
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
	functionName := getCallingFunctionName()
	typeOfCaller := reflect.TypeOf(caller)
	method, _ := typeOfCaller.MethodByName(functionName)
	return s.MethodCalled(method, arguments...)
}

func (s *Stub) MethodCalled(method reflect.Method, arguments ...interface{}) Arguments {
	numberOfReturnArguments := method.Type.NumOut()
	if numberOfReturnArguments == 0 {
		return nil
	}

	foundCall := s.findPredefinedCall(method.Name, arguments...)
	if foundCall == nil && numberOfReturnArguments > 0 {
		errorMessage := fmt.Sprintf("I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", method.Name)
		panic(errorMessage)
	}

	return foundCall.ReturnArguments
}

func (s *Stub) findPredefinedCall(methodName string, arguments ...interface{}) *Call {

	for _, registeredCall := range s.PredefinedCalls {
		if methodName == registeredCall.MethodName {
			for i, argument := range arguments {
				if registeredCall.Arguments[i] != argument {
					return nil
				}
			}
			return registeredCall
		}
	}
	return nil
}

type Double interface {
	On(methodName string, arguments ...interface{}) *Call
	Called(caller interface{}, arguments ...interface{}) Arguments
	MethodCalled(method reflect.Method, arguments ...interface{}) Arguments
}
