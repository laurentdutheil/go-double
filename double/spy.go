package double

import "reflect"

type Spy struct {
	Stub
	ActualCalls []Call
}

func (s *Spy) Called(caller interface{}, arguments ...interface{}) Arguments {
	functionName := getCallingFunctionName()
	typeOfCaller := reflect.TypeOf(caller)
	method, _ := typeOfCaller.MethodByName(functionName)
	return s.MethodCalled(method, arguments...)
}

func (s *Spy) MethodCalled(method reflect.Method, arguments ...interface{}) Arguments {
	call := *NewCall(method.Name, arguments...)
	s.ActualCalls = append(s.ActualCalls, call)

	return s.Stub.MethodCalled(method, arguments...)
}
