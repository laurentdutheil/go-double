package double

import "fmt"

type Stub struct {
	PredefinedCalls []*Call
}

func (s *Stub) On(methodName string, arguments ...interface{}) *Call {
	call := NewCall(methodName, arguments...)
	s.PredefinedCalls = append(s.PredefinedCalls, call)
	return call
}

func (s *Stub) Called(arguments ...interface{}) Arguments {
	functionName := getCallingFunctionName()
	return s.MethodCalled(functionName, arguments...)
}

func (s *Stub) MethodCalled(methodName string, arguments ...interface{}) Arguments {
	foundCall := s.findPredefinedCall(methodName, arguments...)
	if foundCall == nil {
		errorMessage := fmt.Sprintf("I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", methodName)
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
