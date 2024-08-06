package double

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
)

type Mock struct {
	Spy
	inOrder *MocksInOrder
}

func (m *Mock) Called(arguments ...interface{}) Arguments {
	methodInformation := m.getMethodInformation()
	return m.MethodCalled(*methodInformation, arguments...)
}

func (m *Mock) MethodCalled(methodInformation MethodInformation, arguments ...interface{}) Arguments {
	m.recordCallInOrder(methodInformation.Name, arguments...)

	return m.Spy.MethodCalled(methodInformation, arguments...)
}

func (m *Mock) AddActualCall(arguments ...interface{}) {
	functionName := GetCallingFunctionName(2)
	m.recordCallInOrder(functionName, arguments...)
	m.Spy.addActualCall(functionName, arguments)
}

func (m *Mock) AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCalls(methodName)

	return assert.Equal(t, expectedCalls, numberOfCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, numberOfCalls))
}

func (m *Mock) AssertNumberOfCallsWithArguments(t TestingT, expectedCalls int, methodName string, arguments ...interface{}) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCallsWithArguments(methodName, arguments...)

	return assert.Equal(t, expectedCalls, numberOfCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, numberOfCalls))
}

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

func (m *Mock) AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCallsWithArguments(methodName, arguments...)

	if numberOfCalls > 0 {
		return assert.Fail(t, "Should not have called with given arguments",
			fmt.Sprintf("Expected %q to not have been called with:\n%v\nbut actually it was.", methodName, arguments))
	}
	return true
}

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

func (m *Mock) InOrder(inOrder *MocksInOrder) {
	m.inOrder = inOrder
}

func (m *Mock) recordCallInOrder(methodName string, arguments ...interface{}) {
	if m.inOrder != nil {
		call := NewActualCall(methodName, arguments...)
		m.inOrder.addCall(call)
	}
}

type IMock interface {
	ISpy
	AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool
	AssertNumberOfCallsWithArguments(t TestingT, expectedCalls int, methodName string, arguments ...interface{}) bool
	AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool
	AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool
	InOrder(inOrder *MocksInOrder)
}

// Check if Mock implements all methods of IMock
var _ IMock = (*Mock)(nil)
