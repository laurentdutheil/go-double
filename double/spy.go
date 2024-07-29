package double

type Spy[T interface{}] struct {
	Stub[T]
	actualCalls []ActualCall
}

func (s *Spy[T]) Called(arguments ...interface{}) Arguments {
	methodInformation := s.getMethodInformation()
	return s.MethodCalled(*methodInformation, arguments...)
}

func (s *Spy[T]) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	call := NewActualCall(methodInformation.Name, arguments...)
	s.actualCalls = append(s.actualCalls, call)

	return s.Stub.MethodCalled(methodInformation, arguments...)
}

func (s *Spy[T]) AddActualCall(arguments ...interface{}) {
	functionName := GetCallingFunctionName(2)
	call := NewActualCall(functionName, arguments...)
	s.actualCalls = append(s.actualCalls, call)
}

func (s *Spy[T]) NumberOfCalls(methodName string) int {
	count := 0
	for _, call := range s.actualCalls {
		if call.MethodName == methodName {
			count++
		}
	}
	return count
}

func (s *Spy[T]) NumberOfCallsWithArguments(methodName string, arguments ...interface{}) int {
	count := 0
	for _, call := range s.actualCalls {
		if call.isEqual(methodName, arguments) {
			count++
		}
	}
	return count
}

func (s *Spy[T]) ActualCalls() []ActualCall {
	return s.actualCalls
}

type ActualCall struct {
	MethodName string
	Arguments  []interface{}
}

func NewActualCall(methodName string, arguments ...interface{}) ActualCall {
	return ActualCall{MethodName: methodName, Arguments: arguments}
}

func (a ActualCall) isEqual(methodName string, arguments Arguments) bool {
	if a.MethodName != methodName {
		return false
	}
	if len(a.Arguments) != len(arguments) {
		return false
	}

	return arguments.Matches(a.Arguments...)
}
