package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestSpy_Called(t *testing.T) {
	t.Run("Panic if do not use the New constructor method", func(t *testing.T) {
		spy := &SpyExample{}

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { spy.Method() })
	})

	t.Run("Panic if do use the New constructor method incorrectly", func(t *testing.T) {
		spy := New[SpyExample](nil)

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { spy.Method() })
	})

	t.Run("FailNow when private method is used with Called", func(t *testing.T) {
		st := &SpiedTestingT{}
		spy := New[SpyExample](st)
		spy.On("privateMethod").Return(nil)

		st.AssertFailNowWasCalled(t, func() {
			_ = spy.privateMethod()
		})

		assert.Equal(t, "couldn't get the caller method information. 'privateMethod' is private or does not exist\n\tUse MethodCalled instead of Called in stub implementation.", st.errorfFormat)
	})

	t.Run("Register actual call", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		spy.Method()

		assert.Len(t, spy.ActualCalls, 1)
		assert.Equal(t, NewActualCall("Method"), spy.ActualCalls[0])
	})

	t.Run("Register actual call with arguments", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		spy.MethodWithArguments(123, "123", 123.0)

		assert.Len(t, spy.ActualCalls, 1)
		assert.Equal(t, NewActualCall("MethodWithArguments", 123, "123", 123.0), spy.ActualCalls[0])
	})
}

func TestSpy_MethodCalled(t *testing.T) {
	t.Run("Panic if do not use the New constructor method", func(t *testing.T) {
		spy := &SpyExample{}

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { _ = spy.privateMethodWithMethodCalled(1) })
	})

	t.Run("Panic if do use the New constructor method incorrectly", func(t *testing.T) {
		spy := New[SpyExample](nil)

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { _ = spy.privateMethodWithMethodCalled(1) })
	})

	t.Run("Use MethodCalled on private method", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)
		expectedErr := fmt.Errorf("stubbed error")
		spy.On("privateMethodWithMethodCalled", 1).Return(expectedErr)

		err := spy.privateMethodWithMethodCalled(1)

		assert.Equal(t, expectedErr, err)
	})

	t.Run("Register actual call on private method", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)
		expectedErr := fmt.Errorf("stubbed error")
		spy.On("privateMethodWithMethodCalled", 1).Return(expectedErr)

		err := spy.privateMethodWithMethodCalled(1)

		assert.Equal(t, expectedErr, err)
		assert.Len(t, spy.ActualCalls, 1)
		assert.Equal(t, NewActualCall("privateMethodWithMethodCalled", 1), spy.ActualCalls[0])
	})
}

func TestSpy_AddActualCall(t *testing.T) {
	t.Run("Register actual call", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		result := spy.methodOnlyAddActualCall(1)

		assert.Equal(t, 123, result)
		assert.Len(t, spy.ActualCalls, 1)
		assert.Equal(t, NewActualCall("methodOnlyAddActualCall", 1), spy.ActualCalls[0])
	})
}

func TestSpy_NumberOfCalls(t *testing.T) {
	t.Run("Zero call", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		assert.Equal(t, 0, spy.NumberOfCalls("Method"))
	})

	t.Run("Several calls", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		spy.Method()
		spy.Method()

		assert.Equal(t, 2, spy.NumberOfCalls("Method"))
	})
}

func TestSpy_NumberOfCallsWithArguments(t *testing.T) {
	t.Run("Zero call", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		assert.Equal(t, 0, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
	})

	t.Run("Several calls", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		spy.MethodWithArguments(1, "2", 3.0)
		spy.MethodWithArguments(1, "2", 3.0)
		spy.MethodWithArguments(1, "2", 3.0)

		assert.Equal(t, 3, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
	})

	t.Run("One call with wrong arguments", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		spy.MethodWithArguments(0, "2", 3.0)

		assert.Equal(t, 0, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
	})

	t.Run("One call with wrong number of arguments", func(t *testing.T) {
		tt := new(testing.T)
		spy := New[SpyExample](tt)

		spy.MethodWithArguments(0, "2", 3.0)

		assert.Equal(t, 0, spy.NumberOfCallsWithArguments("MethodWithArguments", 1))
	})
}
