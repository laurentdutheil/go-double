package double

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
)

// Mock is a Spy that can do assertion on expected calls
// For an example of its usage, refer to the "Example Usage" section at the top
// of this document.
type Mock struct {
	Spy
	inOrderValidator *InOrderValidator
}

// Called tells the mock object that a method has been called, and gets an array
// of arguments to return.  Fail the test if the call is unexpected (i.e. not preceded by
// appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (m *Mock) Called(arguments ...interface{}) Arguments {
	methodInformation := m.getMethodInformation()
	return m.MethodCalled(*methodInformation, arguments...)
}

// MethodCalled tells the mock object that a method has been called, and gets an array
// of arguments to return.  Fail the test if the call is unexpected (i.e. not preceded by
// appropriate .On .Return() calls)
// If Call.WaitFor is set, blocks until the channel is closed or receives a message.
func (m *Mock) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	m.recordCallInOrder(methodInformation.Name, arguments...)
	return m.Spy.MethodCalled(methodInformation, arguments...)
}

// AddActualCall records the actual call
func (m *Mock) AddActualCall(arguments ...interface{}) {
	functionName := GetCallingFunctionName(2)
	m.recordCallInOrder(functionName, arguments...)
	m.Spy.addActualCall(functionName, arguments)
}

// AssertNumberOfCalls asserts that the method was called expectedCalls times.
func (m *Mock) AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCalls(methodName)

	return assert.Equal(t, expectedCalls, numberOfCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, numberOfCalls))
}

// AssertNumberOfCallsWithArguments asserts that the method was called expectedCalls times.
func (m *Mock) AssertNumberOfCallsWithArguments(t TestingT, expectedCalls int, methodName string, arguments ...interface{}) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCallsWithArguments(methodName, arguments...)

	return assert.Equal(t, expectedCalls, numberOfCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, numberOfCalls))
}

// AssertCalled asserts that the method was called.
// It can produce a false result when an argument is a pointer type and the underlying value changed after calling the mocked method.
func (m *Mock) AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCallsWithArguments(methodName, arguments...)
	if numberOfCalls == 0 {
		var calledWithArgs []string
		for _, call := range m.ActualCalls() {
			if call.MethodName == methodName {
				calledWithArgs = append(calledWithArgs, fmt.Sprintf("%v", call.Arguments))
			}
		}

		if len(calledWithArgs) == 0 {
			return assert.Fail(t, "Should have called with given arguments",
				fmt.Sprintf("Expected %q to have been called with:\n%v\nbut no actual calls happened", methodName, arguments))
		}

		return assert.Fail(t, "Should have called with given arguments", fmt.Sprintf("Expected %q to have been called with:\n%v\nbut actual calls were:\n        %v", methodName, arguments, strings.Join(calledWithArgs, "\n")))
	}

	return true
}

// AssertNotCalled asserts that the method was not called.
// It can produce a false result when an argument is a pointer type and the underlying value changed after calling the mocked method.
func (m *Mock) AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCallsWithArguments(methodName, arguments...)

	if numberOfCalls > 0 {
		return assert.Fail(t, "Should not have called with given arguments",
			fmt.Sprintf("Expected %q to not have been called with:\n%v\nbut actually it was.", methodName, arguments))
	}
	return true
}

// AssertExpectations asserts that everything specified with On and Return was
// in fact called as expected.  Calls may have occurred in any order.
// Deprecated: to respect the 'Arrange, Act, Assert' pattern, consider using the Assert* methods instead
func (m *Mock) AssertExpectations(t TestingT) bool {
	t.Helper()

	result := true
	for _, call := range m.PredefinedCalls() {
		expected := m.AssertCalled(t, call.MethodName, call.Arguments...)
		if expected && !call.calledPredefinedTimes() {
			expected = assert.Fail(t, "Should have called with given arguments",
				fmt.Sprintf("Expected %q to have been called %d times with:\n%v\nbut actually it was called %d times.", call.MethodName, call.times, call.Arguments, call.totalCalls))
		}

		result = result && expected
	}

	return result
}

func (m *Mock) inOrder(inOrderValidator *InOrderValidator) {
	m.inOrderValidator = inOrderValidator
}

func (m *Mock) recordCallInOrder(methodName string, arguments ...interface{}) {
	if m.inOrderValidator != nil {
		call := NewActualCall(methodName, arguments...)
		m.inOrderValidator.addCall(call)
	}
}

type IMock interface {
	ISpy
	AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool
	AssertNumberOfCallsWithArguments(t TestingT, expectedCalls int, methodName string, arguments ...interface{}) bool
	AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool
	AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool
	inOrder(inOrder *InOrderValidator)
}

// Check if Mock implements all methods of IMock
var _ IMock = (*Mock)(nil)
