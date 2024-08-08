package double

type Spy struct {
	Stub
	actualCalls []ActualCall
}

func (s *Spy) Called(arguments ...interface{}) Arguments {
	methodInformation := s.getMethodInformation()
	return s.MethodCalled(*methodInformation, arguments...)
}

func (s *Spy) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	s.addActualCall(methodInformation.Name, arguments)
	return s.Stub.MethodCalled(methodInformation, arguments...)
}

func (s *Spy) AddActualCall(arguments ...interface{}) {
	functionName := GetCallingFunctionName(2)
	s.addActualCall(functionName, arguments)
}

func (s *Spy) addActualCall(methodName string, arguments []interface{}) {
	call := NewActualCall(methodName, arguments...)
	s.actualCalls = append(s.actualCalls, call)
}

func (s *Spy) NumberOfCalls(methodName string) int {
	count := 0
	for _, call := range s.actualCalls {
		if call.MethodName == methodName {
			count++
		}
	}
	return count
}

func (s *Spy) NumberOfCallsWithArguments(methodName string, arguments ...interface{}) int {
	count := 0
	for _, call := range s.actualCalls {
		if call.matches(methodName, arguments) {
			count++
		}
	}
	return count
}

func (s *Spy) ActualCalls() []ActualCall {
	return s.actualCalls
}

type ActualCall struct {
	MethodName string
	Arguments  []interface{}
}

func NewActualCall(methodName string, arguments ...interface{}) ActualCall {
	return ActualCall{MethodName: methodName, Arguments: arguments}
}

func (a ActualCall) matches(methodName string, arguments Arguments) bool {
	if a.MethodName != methodName {
		return false
	}

	return arguments.Matches(a.Arguments...)
}

type ISpy interface {
	IStub
	AddActualCall(arguments ...interface{})
	NumberOfCalls(methodName string) int
	NumberOfCallsWithArguments(methodName string, arguments ...interface{}) int
	ActualCalls() []ActualCall
}

// Check if Spy implements all methods of ISpy
var _ ISpy = (*Spy)(nil)
