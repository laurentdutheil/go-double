package double

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
)

type Mock[T interface{}] struct {
	Spy[T]
}

func (m *Mock[T]) AssertNumberOfCalls(t TestingT, methodName string, expectedCalls int) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCalls(methodName)

	return assert.Equal(t, expectedCalls, numberOfCalls, fmt.Sprintf("Expected number of calls (%d) does not match the actual number of calls (%d).", expectedCalls, numberOfCalls))
}

func (m *Mock[T]) AssertCalled(t TestingT, methodName string, arguments ...interface{}) bool {
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

func (m *Mock[T]) AssertNotCalled(t TestingT, methodName string, arguments ...interface{}) bool {
	t.Helper()

	numberOfCalls := m.NumberOfCallsWithArguments(methodName, arguments...)

	if numberOfCalls > 0 {
		return assert.Fail(t, "Should not have called with given arguments",
			fmt.Sprintf("Expected %q to not have been called with:\n%v\nbut actually it was.", methodName, arguments))
	}
	return true
}

func (m *Mock[T]) AssertExpectations(t TestingT) bool {
	t.Helper()

	result := true
	for _, call := range m.PredefinedCalls() {
		assertCalled := m.AssertCalled(t, call.MethodName, call.Arguments...)
		if assertCalled {
			if !call.calledPredefinedTimes() {
				result = false
				assert.Fail(t, "Should have called with given arguments",
					fmt.Sprintf("Expected %q to have been called %d times with:\n%v\nbut actually it was called %d times.", call.MethodName, call.times, call.Arguments, call.totalCalls))
			}
		}

		result = result && assertCalled
	}

	return result
}
