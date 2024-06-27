package double

import (
	"runtime"
	"strings"
)

type Stub struct {
	RegisteredCalls []*Call
	Calls           []Call
}

func (s *Stub) On(methodName string, args ...interface{}) *Call {
	call := NewCall(methodName, args...)
	s.RegisteredCalls = append(s.RegisteredCalls, call)
	return call
}

func (s *Stub) Called(args ...interface{}) Arguments {
	functionName := s.getCallingFunctionName()
	call := *NewCall(functionName, args...)
	s.Calls = append(s.Calls, call)

	foundCall := s.findRegisteredCall(call.MethodName)
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

func (s *Stub) findRegisteredCall(methodName string) *Call {
	for _, registeredCall := range s.RegisteredCalls {
		if methodName == registeredCall.MethodName {
			return registeredCall
		}
	}
	return nil
}
