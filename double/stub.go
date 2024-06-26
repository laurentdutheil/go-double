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

func (s *Stub) Called(args ...interface{}) {
	functionName := s.getCallingFunctionName()
	s.Calls = append(s.Calls, *NewCall(functionName, args...))
}

func (s *Stub) getCallingFunctionName() string {
	pc, _, _, _ := runtime.Caller(2)
	functionPath := runtime.FuncForPC(pc).Name()
	parts := strings.Split(functionPath, ".")
	return parts[len(parts)-1]
}
