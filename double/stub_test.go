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

	aInt, err := stub.MethodWithReturnArguments()

	assert.Equal(t, expectedInt, aInt)
	assert.Equal(t, expectedErr, err)
}

func TestStub_ReturnPredefinedReturnArgumentsWithArgumentsChecking(t *testing.T) {
	stub := &StubExample{}
	expectedInt := 1
	expectedErr := fmt.Errorf("stubbed error")
	stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)

	aInt, err := stub.MethodWithArgumentsAndReturnArguments(123, "123", 123.0)

	assert.Equal(t, expectedInt, aInt)
	assert.Equal(t, expectedErr, err)
}

func TestStub_PanicsWithArgumentsChecking(t *testing.T) {
	stub := &StubExample{}
	expectedInt := 1
	expectedErr := fmt.Errorf("stubbed error")
	stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)

	expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithArgumentsAndReturnArguments\").Return(...) first"
	assert.PanicsWithValue(t, expectedError, func() { _, _ = stub.MethodWithArgumentsAndReturnArguments(12, "", 1.0) })
}

func TestStub_Called_ShouldNotPanicOnMethodWithoutReturnArgument(t *testing.T) {
	stub := &StubExample{}

	assert.NotPanics(t, func() { stub.Method() })
}

func TestStub_ReturnPredefinedReturnArgumentsOnce(t *testing.T) {
	stub := &StubExample{}
	expectedInt := 1
	expectedErr := fmt.Errorf("stubbed error")
	stub.On("MethodWithReturnArguments").Return(expectedInt, expectedErr).Once()

	aInt, err := stub.MethodWithReturnArguments()
	assert.Equal(t, expectedInt, aInt)
	assert.Equal(t, expectedErr, err)

	expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithReturnArguments\").Return(...) first"
	assert.PanicsWithValue(t, expectedError, func() { _, _ = stub.MethodWithReturnArguments() })
}

func TestStub_ReturnPredefinedReturnDifferentArgumentsOnDifferentCall(t *testing.T) {
	stub := &StubExample{}
	expectedErr := fmt.Errorf("stubbed error")
	stub.On("MethodWithReturnArguments").Return(1, expectedErr).Once()
	stub.On("MethodWithReturnArguments").Return(2, expectedErr).Once()

	aInt1, err1 := stub.MethodWithReturnArguments()
	assert.Equal(t, 1, aInt1)
	assert.Equal(t, expectedErr, err1)

	aInt2, err2 := stub.MethodWithReturnArguments()
	assert.Equal(t, 2, aInt2)
	assert.Equal(t, expectedErr, err2)

	expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithReturnArguments\").Return(...) first"
	assert.PanicsWithValue(t, expectedError, func() { _, _ = stub.MethodWithReturnArguments() })
}
