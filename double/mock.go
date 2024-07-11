package double

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type Mock struct {
	Spy
}

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
	Helper()
}

// Check if TestingT interface can wrap testing.T
var _ TestingT = (*testing.T)(nil)

func (m *Mock) AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCalls(methodName)

	return assert.Equal(t, expectedCalls, numberOfCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, numberOfCalls))
}

func (m *Mock) AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCallsWithArguments(methodName, arguments...)
	if numberOfCalls == 0 {
		var calledWithArgs []string
		for _, call := range m.ActualCalls {
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
