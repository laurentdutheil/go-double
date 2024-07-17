package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	. "github.com/laurentdutheil/go-double/double"
)

func TestStub_On(t *testing.T) {
	t.Run("Predefine method name", func(t *testing.T) {
		tt := new(testing.T)
		stub := New[StubExample](tt)

		call := stub.On("Method")

		assert.Equal(t, "Method", call.MethodName)
		assert.Contains(t, stub.PredefinedCalls(), call)
	})

	t.Run("Predefine method name and arguments", func(t *testing.T) {
		tt := new(testing.T)
		stub := New[StubExample](tt)

		call := stub.On("MethodWithArguments", 1, "2", 3.0)

		assert.Equal(t, "MethodWithArguments", call.MethodName)
		assert.Contains(t, stub.PredefinedCalls(), call)
		assert.Len(t, call.Arguments, 3)
		assert.Contains(t, call.Arguments, 1)
		assert.Contains(t, call.Arguments, "2")
		assert.Contains(t, call.Arguments, 3.0)
	})
}

func TestStub_Called(t *testing.T) {
	t.Run("Panic if do not use the New constructor method", func(t *testing.T) {
		stub := &StubExample{}

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { stub.Method() })
	})

	t.Run("Panic if do use the New constructor method incorrectly", func(t *testing.T) {
		stub := New[StubExample](nil)

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { stub.Method() })
	})
}

func TestStub_On_Return(t *testing.T) {
	t.Run("Predefine return arguments", func(t *testing.T) {
		tt := new(testing.T)
		stub := New[StubExample](tt)
		expectedInt := 1
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(expectedInt, expectedErr)

		aInt, err := stub.MethodWithReturnArguments()

		assert.Equal(t, expectedInt, aInt)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("Predefine return arguments with arguments checking", func(t *testing.T) {
		tt := new(testing.T)
		stub := New[StubExample](tt)
		expectedInt := 1
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)

		aInt, err := stub.MethodWithArgumentsAndReturnArguments(123, "123", 123.0)

		assert.Equal(t, expectedInt, aInt)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("FailNow when arguments don't match", func(t *testing.T) {
		st := &SpiedTestingT{}
		stub := New[StubExample](st)
		expectedInt := 1
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)

		st.AssertFailNowWasCalled(t, func() {
			_, _ = stub.MethodWithArgumentsAndReturnArguments(12, "", 1.0)
		})
		assert.Equal(t, "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"%s\").Return(...) first", st.errorfFormat)
		errorMethodeName := st.errorfArgs[0]
		assert.Equal(t, errorMethodeName, "MethodWithArgumentsAndReturnArguments")
	})

	t.Run("Don't panic when method have no return arguments. Even if there is no predefined call", func(t *testing.T) {
		tt := new(testing.T)
		stub := New[StubExample](tt)

		assert.NotPanics(t, func() { stub.Method() })
	})
}

func TestStub_Times(t *testing.T) {
	t.Run("Return predefined return arguments once. And FailNow on the additional call", func(t *testing.T) {
		st := &SpiedTestingT{}
		stub := New[StubExample](st)
		expectedInt := 1
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(expectedInt, expectedErr).Once()

		aInt, err := stub.MethodWithReturnArguments()
		assert.Equal(t, expectedInt, aInt)
		assert.Equal(t, expectedErr, err)

		st.AssertFailNowWasCalled(t, func() {
			_, _ = stub.MethodWithReturnArguments()
		})
	})

	t.Run("Return predefined return arguments twice. And panic on the additional call", func(t *testing.T) {
		st := &SpiedTestingT{}
		stub := New[StubExample](st)
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(1, expectedErr).Twice()

		for i := 1; i <= 2; i++ {
			aInt, err := stub.MethodWithReturnArguments()
			assert.Equal(t, 1, aInt)
			assert.Equal(t, expectedErr, err)
		}

		st.AssertFailNowWasCalled(t, func() {
			_, _ = stub.MethodWithReturnArguments()
		})
	})

	t.Run("Return predefined return arguments n times. And panic on the additional call", func(t *testing.T) {
		st := &SpiedTestingT{}
		stub := New[StubExample](st)
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(1, expectedErr).Times(4)

		for i := 1; i <= 4; i++ {
			aInt, err := stub.MethodWithReturnArguments()
			assert.Equal(t, 1, aInt)
			assert.Equal(t, expectedErr, err)
		}

		st.AssertFailNowWasCalled(t, func() {
			_, _ = stub.MethodWithReturnArguments()
		})
	})

	t.Run("Return different predefined return arguments. And panic on the additional call", func(t *testing.T) {
		st := &SpiedTestingT{}
		stub := New[StubExample](st)
		expectedErr := fmt.Errorf("stubbed error")
		stub.On("MethodWithReturnArguments").Return(1, expectedErr).Once()
		stub.On("MethodWithReturnArguments").Return(2, expectedErr).Once()

		for i := 1; i <= 2; i++ {
			aInt, err := stub.MethodWithReturnArguments()
			assert.Equal(t, i, aInt)
			assert.Equal(t, expectedErr, err)
		}

		st.AssertFailNowWasCalled(t, func() {
			_, _ = stub.MethodWithReturnArguments()
		})
	})
}

func TestStub_On_Panic(t *testing.T) {
	t.Run("Panic when predefined call say so", func(t *testing.T) {
		tt := new(testing.T)
		stub := New[StubExample](tt)
		stub.On("Method").Panic("panic message for example method")

		assert.PanicsWithValue(t, "panic message for example method", func() {
			stub.Method()
		})
	})
}

func TestStub_On_WailUntil(t *testing.T) {
	t.Run("Wait until the duration", func(t *testing.T) {
		tt := new(testing.T)
		stub := New[StubExample](tt)
		stub.On("Method").WaitUntil(time.After(10 * time.Millisecond))

		done := make(chan string)
		go func() {
			stub.Method()
			done <- "done"
		}()

		// check that it is not done before
		select {
		case <-time.After(5 * time.Millisecond):
			// Pass
		case <-done:
			assert.Fail(t, "Have to wait until the duration")
		}

		// check that it is done after
		select {
		case <-time.After(20 * time.Millisecond):
			assert.Fail(t, "The wait is too long")
		case msg := <-done:
			assert.Equal(t, "done", msg)
		}

	})
}

func TestStub_On_After(t *testing.T) {
	t.Run("Wait until the duration", func(t *testing.T) {
		tt := new(testing.T)
		stub := New[StubExample](tt)
		stub.On("Method").After(10 * time.Millisecond)

		done := make(chan string)
		go func() {
			stub.Method()
			done <- "done"
		}()

		// check that it is not done before
		select {
		case <-time.After(5 * time.Millisecond):
			// Pass
		case <-done:
			assert.Fail(t, "Have to wait until the duration")
		}

		// check that it is done after
		select {
		case <-time.After(20 * time.Millisecond):
			assert.Fail(t, "The wait is too long")
		case msg := <-done:
			assert.Equal(t, "done", msg)
		}

	})
}
