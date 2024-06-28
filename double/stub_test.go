package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
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

func TestStubOn_PredefineMethodNameWithReturnArguments(t *testing.T) {
	stub := &StubExample{}

	call := stub.On("Method").Return(1, nil)

	assert.Equal(t, "Method", call.MethodName)
	assert.Contains(t, stub.PredefinedCalls, call)
	assert.Contains(t, call.ReturnArguments, 1)
	assert.Contains(t, call.ReturnArguments, nil)
}

func TestStub_CallIsPredefined(t *testing.T) {
	stub := &StubExample{}
	stub.On("Method")
	sut := &SUTExample{stub}

	sut.method()

	assert.Len(t, stub.ActualCalls, 1)
	assert.Equal(t, *NewCall("Method"), stub.ActualCalls[0])
}

func TestStub_CallWithArgumentsIsPredefined(t *testing.T) {
	stub := &StubExample{}
	stub.On("MethodWithArguments", 123, "123", 123.0)
	sut := &SUTExample{stub}

	sut.methodWithArguments(123)

	assert.Len(t, stub.ActualCalls, 1)
	assert.Equal(t, *NewCall("MethodWithArguments", 123, "123", 123.0), stub.ActualCalls[0])
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
	assert.Len(t, stub.ActualCalls, 1)
	assert.Equal(t, *NewCall("MethodWithReturnArguments"), stub.ActualCalls[0])
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
	assert.Len(t, stub.ActualCalls, 1)
	assert.Equal(t, *NewCall("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0), stub.ActualCalls[0])
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
