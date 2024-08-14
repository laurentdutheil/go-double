package double

// Spy is a Stub that record actual calls
// For an example of its usage, refer to the "Example Usage" section at the top
// of this document.
type Spy struct {
	Stub
	actualCalls ActualCalls
}

// Called tells the spy object that a method has been called, and gets an array
// of arguments to return.  Fail the test if the call is unexpected (i.e. not preceded by
// appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (s *Spy) Called(arguments ...interface{}) Arguments {
	methodInformation := s.getMethodInformation()
	return s.MethodCalled(*methodInformation, arguments...)
}

// MethodCalled tells the spy object that a method has been called, and gets an array
// of arguments to return.  Fail the test if the call is unexpected (i.e. not preceded by
// appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (s *Spy) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	s.actualCalls.append(methodInformation.Name, arguments)
	return s.Stub.MethodCalled(methodInformation, arguments...)
}

// AddActualCall records the actual call
func (s *Spy) AddActualCall(arguments ...interface{}) {
	functionName := GetCallingFunctionName(2)
	s.actualCalls.append(functionName, arguments)
}

// NumberOfCalls return the number of calls of the method name passed in parameter
func (s *Spy) NumberOfCalls(methodName string) int {
	predicate := func(call ActualCall) bool {
		return call.MethodName == methodName
	}
	return s.actualCalls.count(predicate)
}

// NumberOfCallsWithArguments return the number of calls of the method with the specified arguments
func (s *Spy) NumberOfCallsWithArguments(methodName string, arguments ...interface{}) int {
	predicate := func(call ActualCall) bool {
		return call.matches(methodName, arguments)
	}
	return s.actualCalls.count(predicate)
}

// ActualCalls return the actual calls recorded by the Spy
func (s *Spy) ActualCalls() []ActualCall {
	return s.actualCalls
}

// ActualCall record the information of an actual call
type ActualCall struct {
	MethodName string
	Arguments  []interface{}
}

// NewActualCall constructor
func NewActualCall(methodName string, arguments ...interface{}) ActualCall {
	return ActualCall{MethodName: methodName, Arguments: arguments}
}

func (a ActualCall) matches(methodName string, arguments Arguments) bool {
	if a.MethodName != methodName {
		return false
	}

	return arguments.Matches(a.Arguments...)
}

type ActualCalls []ActualCall

func (c *ActualCalls) append(methodName string, arguments []interface{}) {
	call := NewActualCall(methodName, arguments...)
	*c = append(*c, call)
}

func (c *ActualCalls) count(predicate func(ActualCall) bool) int {
	count := 0
	for _, call := range *c {
		if predicate(call) {
			count++
		}
	}
	return count
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
