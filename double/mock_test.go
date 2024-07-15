package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestMock_Called(t *testing.T) {
	t.Run("Panic if do not use the New constructor method", func(t *testing.T) {
		mock := &MockExample{}

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { mock.Method() })
	})

	t.Run("Panic if do use the New constructor method incorrectly", func(t *testing.T) {
		mock := New[StubExample](nil)

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { mock.Method() })
	})
}

func TestMock_AssertNumberOfCalls(t *testing.T) {

	t.Run("t.Helper is called", func(t *testing.T) {
		st := &SpiedTestingT{}
		mock := New[MockExample](st)

		mock.AssertNumberOfCalls(st, "Method", 1)

		assert.True(t, st.helperCalled)
	})

	t.Run("Return false when number of calls is incorrect", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)

		result := mock.AssertNumberOfCalls(tt, "Method", 1)

		assert.False(t, result)
	})

	t.Run("Return true when number of calls is correct", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)
		mock.Method()

		result := mock.AssertNumberOfCalls(tt, "Method", 1)

		assert.True(t, result)
	})

	t.Run("t.Errorf is called with right message when number of calls is incorrect", func(t *testing.T) {
		st := &SpiedTestingT{}
		mock := New[MockExample](st)
		mock.Method()

		mock.AssertNumberOfCalls(st, "Method", 2)

		assert.Equal(t, "\n%s", st.errorfFormat)
		errorMessage := st.errorfArgs[0]
		assert.Contains(t, errorMessage, "Error Trace:")
		assert.Contains(t, errorMessage, "Expected number of calls (2) does not match the actual number of calls (1).")
	})

}

func TestMock_AssertCall(t *testing.T) {

	t.Run("t.Helper is called", func(t *testing.T) {
		st := &SpiedTestingT{}
		mock := New[MockExample](st)

		mock.AssertCalled(st, "Method")

		assert.True(t, st.helperCalled)
	})

	t.Run("Return false when method is not called", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)

		result := mock.AssertCalled(tt, "Method")

		assert.False(t, result)
	})

	t.Run("Return true when method is called", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)
		mock.Method()

		result := mock.AssertCalled(tt, "Method")

		assert.True(t, result)
	})

	t.Run("t.Errorf is called with right message when method is not called", func(t *testing.T) {
		st := &SpiedTestingT{}
		mock := New[MockExample](st)

		mock.AssertCalled(st, "Method")

		assert.Equal(t, "\n%s", st.errorfFormat)
		errorMessage := st.errorfArgs[0]
		assert.Contains(t, errorMessage, "Error Trace:")
		assert.Contains(t, errorMessage, "Should have called with given arguments")
		assert.Contains(t, errorMessage, "Expected \"Method\" to have been called with:")
		assert.Contains(t, errorMessage, "[]")
		assert.Contains(t, errorMessage, "but no actual calls happened")
	})

	t.Run("Return true when method is called with right arguments", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)
		mock.MethodWithArguments(1, "1", 1.0)

		result := mock.AssertCalled(tt, "MethodWithArguments", 1, "1", 1.0)

		assert.True(t, result)
	})

	t.Run("Return false when method is called with wrong arguments", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)
		mock.MethodWithArguments(2, "1", 1.0)

		result := mock.AssertCalled(tt, "MethodWithArguments", 1, "1", 1.0)

		assert.False(t, result)
	})

	t.Run("t.Errorf is called with right message when method is called with wrong arguments", func(t *testing.T) {
		st := &SpiedTestingT{}
		mock := New[MockExample](st)
		mock.MethodWithArguments(2, "1", 1.0)
		mock.MethodWithArguments(1, "3", 1.2)
		mock.MethodWithOneArgument(4)

		mock.AssertCalled(st, "MethodWithArguments", 1, "1", 1.0)

		assert.Equal(t, "\n%s", st.errorfFormat)
		errorMessage := st.errorfArgs[0]
		assert.Contains(t, errorMessage, "Error Trace:")
		assert.Contains(t, errorMessage, "Should have called with given arguments")
		assert.Contains(t, errorMessage, "Expected \"MethodWithArguments\" to have been called with:")
		assert.Contains(t, errorMessage, "[1 1 1]")
		assert.Contains(t, errorMessage, "but actual calls were:")
		assert.Contains(t, errorMessage, "[2 1 1]")
		assert.Contains(t, errorMessage, "[1 3 1.2]")
		assert.NotContains(t, errorMessage, "[4]")
	})
}

func TestMock_AssertNotCall(t *testing.T) {

	t.Run("t.Helper is called", func(t *testing.T) {
		st := &SpiedTestingT{}
		mock := New[MockExample](st)

		mock.AssertNotCalled(st, "Method")

		assert.True(t, st.helperCalled)
	})

	t.Run("Return false when method is called", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)
		mock.Method()

		result := mock.AssertNotCalled(tt, "Method")

		assert.False(t, result)
	})

	t.Run("Return true when method is not called", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)

		result := mock.AssertNotCalled(tt, "Method")

		assert.True(t, result)
	})

	t.Run("t.Errorf is called with right message when method is called", func(t *testing.T) {
		st := &SpiedTestingT{}
		mock := New[MockExample](st)
		mock.Method()

		mock.AssertNotCalled(st, "Method")

		assert.Equal(t, "\n%s", st.errorfFormat)
		errorMessage := st.errorfArgs[0]
		assert.Contains(t, errorMessage, "Error Trace:")
		assert.Contains(t, errorMessage, "Should not have called with given arguments")
		assert.Contains(t, errorMessage, "Expected \"Method\" to not have been called with:")
		assert.Contains(t, errorMessage, "[]")
		assert.Contains(t, errorMessage, "but actually it was.")
	})

	t.Run("Return true when method is called with other arguments", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)
		mock.MethodWithArguments(2, "1", 1.0)

		result := mock.AssertNotCalled(tt, "MethodWithArguments", 1, "1", 1.0)

		assert.True(t, result)
	})

	t.Run("Return false when method is called with same arguments", func(t *testing.T) {
		tt := new(testing.T)
		mock := New[MockExample](tt)
		mock.MethodWithArguments(1, "1", 1.0)

		result := mock.AssertNotCalled(tt, "MethodWithArguments", 1, "1", 1.0)

		assert.False(t, result)
	})

	t.Run("t.Errorf is called with right message when method is called with same arguments", func(t *testing.T) {
		st := &SpiedTestingT{}
		mock := New[MockExample](st)
		mock.MethodWithArguments(1, "3", 1.2)

		mock.AssertNotCalled(st, "MethodWithArguments", 1, "3", 1.2)

		assert.Equal(t, "\n%s", st.errorfFormat)
		errorMessage := st.errorfArgs[0]
		assert.Contains(t, errorMessage, "Error Trace:")
		assert.Contains(t, errorMessage, "Should not have called with given arguments")
		assert.Contains(t, errorMessage, "Expected \"MethodWithArguments\" to not have been called with:")
		assert.Contains(t, errorMessage, "[1 3 1.2]")
		assert.Contains(t, errorMessage, "but actually it was.")
	})
}
