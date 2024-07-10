package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStub_On(t *testing.T) {
	t.Run("Predefine method name", func(t *testing.T) {
		stub := &StubExample{}

		call := stub.On("Method")

		assert.Equal(t, "Method", call.MethodName)
		assert.Contains(t, stub.PredefinedCalls, call)
	})

	t.Run("Predefine method name and arguments", func(t *testing.T) {
		stub := &StubExample{}

		call := stub.On("MethodWithArguments", 1, "2", 3.0)

		assert.Equal(t, "MethodWithArguments", call.MethodName)
		assert.Contains(t, stub.PredefinedCalls, call)
		assert.Len(t, call.Arguments, 3)
		assert.Contains(t, call.Arguments, 1)
		assert.Contains(t, call.Arguments, "2")
		assert.Contains(t, call.Arguments, 3.0)
	})
}

func TestStub_On_Return(t *testing.T) {
	t.Run("Predefine return arguments", func(t *testing.T) {
		stub := &StubExample{}
		expectedInt := 1
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(expectedInt, expectedErr)

		aInt, err := stub.MethodWithReturnArguments()

		assert.Equal(t, expectedInt, aInt)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Predefine return arguments with arguments checking", func(t *testing.T) {
		stub := &StubExample{}
		expectedInt := 1
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)

		aInt, err := stub.MethodWithArgumentsAndReturnArguments(123, "123", 123.0)

		assert.Equal(t, expectedInt, aInt)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Panics when arguments don't match", func(t *testing.T) {
		stub := &StubExample{}
		expectedInt := 1
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)

		expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithArgumentsAndReturnArguments\").Return(...) first"
		assert.PanicsWithValue(t, expectedError, func() { _, _ = stub.MethodWithArgumentsAndReturnArguments(12, "", 1.0) })
	})

	t.Run("Don't panic when method have no return arguments. Even if there is no predefined call", func(t *testing.T) {
		stub := &StubExample{}

		assert.NotPanics(t, func() { stub.Method() })
	})
}

func TestStub_Times(t *testing.T) {
	t.Run("Return predefined return arguments once. And panic on the additional call", func(t *testing.T) {
		stub := &StubExample{}
		expectedInt := 1
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(expectedInt, expectedErr).Once()

		aInt, err := stub.MethodWithReturnArguments()
		assert.Equal(t, expectedInt, aInt)
		assert.Equal(t, expectedErr, err)

		expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithReturnArguments\").Return(...) first"
		assert.PanicsWithValue(t, expectedError, func() { _, _ = stub.MethodWithReturnArguments() })
	})

	t.Run("Return predefined return arguments twice. And panic on the additional call", func(t *testing.T) {
		stub := &StubExample{}
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(1, expectedErr).Twice()

		for i := 1; i <= 2; i++ {
			aInt, err := stub.MethodWithReturnArguments()
			assert.Equal(t, 1, aInt)
			assert.Equal(t, expectedErr, err)
		}

		expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithReturnArguments\").Return(...) first"
		assert.PanicsWithValue(t, expectedError, func() { _, _ = stub.MethodWithReturnArguments() })
	})

	t.Run("Return predefined return arguments n times. And panic on the additional call", func(t *testing.T) {
		stub := &StubExample{}
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(1, expectedErr).Times(4)

		for i := 1; i <= 4; i++ {
			aInt, err := stub.MethodWithReturnArguments()
			assert.Equal(t, 1, aInt)
			assert.Equal(t, expectedErr, err)
		}

		expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithReturnArguments\").Return(...) first"
		assert.PanicsWithValue(t, expectedError, func() { _, _ = stub.MethodWithReturnArguments() })
	})

	t.Run("Return different predefined return arguments. And panic on the additional call", func(t *testing.T) {
		stub := &StubExample{}
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(1, expectedErr).Once()
		stub.On("MethodWithReturnArguments").Return(2, expectedErr).Once()

		for i := 1; i <= 2; i++ {
			aInt, err := stub.MethodWithReturnArguments()
			assert.Equal(t, i, aInt)
			assert.Equal(t, expectedErr, err)
		}

		expectedError := "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithReturnArguments\").Return(...) first"
		assert.PanicsWithValue(t, expectedError, func() { _, _ = stub.MethodWithReturnArguments() })
	})
}
