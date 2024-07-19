package double

import "reflect"

type Spy[T interface{}] struct {
	Stub[T]
	ActualCalls []ActualCall
}

func (s *Spy[T]) Called(arguments ...interface{}) Arguments {
	if s.t == nil {
		panic("Please use double.New constructor to initialize correctly.")
	}

	methodInformation := GetCallingMethodInformation(s.caller)
	return s.MethodCalled(methodInformation, arguments...)
}

func (s *Spy[T]) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	call := NewActualCall(methodInformation.Name, arguments...)
	s.ActualCalls = append(s.ActualCalls, call)

	return s.Stub.MethodCalled(methodInformation, arguments...)
}

func (s *Spy[T]) NumberOfCalls(methodName string) int {
	count := 0
	for _, call := range s.ActualCalls {
		if call.MethodName == methodName {
			count++
		}
	}
	return count
}

func (s *Spy[T]) NumberOfCallsWithArguments(methodName string, arguments ...interface{}) int {
	count := 0
	for _, call := range s.ActualCalls {
		if call.isEqual(methodName, arguments) {
			count++
		}
	}
	return count
}

type ActualCall struct {
	MethodName string
	Arguments  []interface{}
}

func NewActualCall(methodName string, arguments ...interface{}) ActualCall {
	return ActualCall{MethodName: methodName, Arguments: arguments}
}

func (a ActualCall) isEqual(methodName string, arguments []interface{}) bool {
	if a.MethodName != methodName {
		return false
	}
	if len(a.Arguments) != len(arguments) {
		return false
	}
	for i, v := range a.Arguments {
		if !reflect.DeepEqual(arguments[i], v) {
			return false
		}
	}
	return true
}
