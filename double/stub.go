package double

import (
	"runtime"
	"strings"
)

type Stub struct {
	PredefinedCalls []*Call
	ActualCalls     []Call
}

func (s *Stub) On(methodName string, arguments ...interface{}) *Call {
	call := NewCall(methodName, arguments...)
	s.PredefinedCalls = append(s.PredefinedCalls, call)
	return call
}

func (s *Stub) Called(arguments ...interface{}) Arguments {
	functionName := s.getCallingFunctionName()
	call := *NewCall(functionName, arguments...)
	s.ActualCalls = append(s.ActualCalls, call)

	foundCall := s.findPredefinedCall(call.MethodName)
	if foundCall == nil {
		return nil
	}

	return foundCall.ReturnArguments
}

func (s *Stub) getCallingFunctionName() string {
	pc, _, _, _ := runtime.Caller(2)
	functionPath := runtime.FuncForPC(pc).Name()
	parts := strings.Split(functionPath, ".")
	return parts[len(parts)-1]
}

func (s *Stub) findPredefinedCall(methodName string) *Call {
	for _, registeredCall := range s.PredefinedCalls {
		if methodName == registeredCall.MethodName {
			return registeredCall
		}
	}
	return nil
}
