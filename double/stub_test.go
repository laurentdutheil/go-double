package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStubOn_PredefineMethodName(t *testing.T) {
	stub := &StubExample{}

	call := stub.On("Method")

	assert.Equal(t, "Method", call.MethodName)
	assert.Contains(t, stub.PredefinedCalls, call)
}

func TestStubOn_PredefineMethodNameAndArguments(t *testing.T) {
	stub := &StubExample{}

	call := stub.On("MethodWithArguments", 1, "2", 3.0)

	assert.Equal(t, "MethodWithArguments", call.MethodName)
	assert.Contains(t, stub.PredefinedCalls, call)
	assert.Len(t, call.Arguments, 3)
	assert.Contains(t, call.Arguments, 1)
	assert.Contains(t, call.Arguments, "2")
	assert.Contains(t, call.Arguments, 3.0)
}

func TestStub_ReturnPredefinedReturnArguments(t *testing.T) {
	stub := &StubExample{}
	expectedInt := 1
	expectedErr := fmt.Errorf("stubbed error")
	stub.On("MethodWithReturnArguments").Return(expectedInt, expectedErr)
	sut := &SUTExample{stub}

	aInt, err := sut.methodWithReturnArguments()

	assert.Equal(t, expectedInt, aInt)
	assert.Equal(t, expectedErr, err)
}

func TestStub_ReturnPredefinedReturnArgumentsWithArgumentsChecking(t *testing.T) {
	stub := &StubExample{}
	expectedInt := 1
	expectedErr := fmt.Errorf("stubbed error")
	stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)
	sut := &SUTExample{stub}

	aInt, err := sut.methodWithArgumentsAndReturnArguments(123)

	assert.Equal(t, expectedInt, aInt)
	assert.Equal(t, expectedErr, err)
}

func TestStub_PanicsWithArgumentsChecking(t *testing.T) {
	stub := &StubExample{}
	expectedInt := 1
	expectedErr := fmt.Errorf("stubbed error")
	stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)
	sut := &SUTExample{stub}

	expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithArgumentsAndReturnArguments\").Return(...) first"
	assert.PanicsWithValue(t, expectedError, func() { _, _ = sut.methodWithArgumentsAndReturnArguments(12) })
}
