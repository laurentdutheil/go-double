package double_test

import (
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
		spy := New[StubExample](nil)

		expectedMessage := "Please use double.New constructor to initialize correctly."
		assert.PanicsWithValue(t, expectedMessage, func() { spy.Method() })
	})
}

func TestSpy_Called_RegisterActualCall(t *testing.T) {
	tt := new(testing.T)
	spy := New[SpyExample](tt)

	spy.Method()

	assert.Len(t, spy.ActualCalls, 1)
	assert.Equal(t, NewActualCall("Method"), spy.ActualCalls[0])
}

func TestSpy_Called_RegisterActualCallWithArguments(t *testing.T) {
	tt := new(testing.T)
	spy := New[SpyExample](tt)

	spy.MethodWithArguments(123, "123", 123.0)

	assert.Len(t, spy.ActualCalls, 1)
	assert.Equal(t, NewActualCall("MethodWithArguments", 123, "123", 123.0), spy.ActualCalls[0])
}

func TestSpy_NumberOfCall_ZeroCall(t *testing.T) {
	tt := new(testing.T)
	spy := New[SpyExample](tt)

	assert.Equal(t, 0, spy.NumberOfCalls("Method"))
}

func TestSpy_NumberOfCall_SeveralCalls(t *testing.T) {
	tt := new(testing.T)
	spy := New[SpyExample](tt)

	spy.Method()
	spy.Method()

	assert.Equal(t, 2, spy.NumberOfCalls("Method"))
}

func TestSpy_NumberOfCallWithArguments_ZeroCall(t *testing.T) {
	tt := new(testing.T)
	spy := New[SpyExample](tt)

	assert.Equal(t, 0, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
}

func TestSpy_NumberOfCallWithArguments_OneCallWithWrongArguments(t *testing.T) {
	tt := new(testing.T)
	spy := New[SpyExample](tt)

	spy.MethodWithArguments(0, "2", 3.0)

	assert.Equal(t, 0, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
}

func TestSpy_NumberOfCallWithArguments_OneCallWithWrongNumberOfArguments(t *testing.T) {
	tt := new(testing.T)
	spy := New[SpyExample](tt)

	spy.MethodWithArguments(0, "2", 3.0)

	assert.Equal(t, 0, spy.NumberOfCallsWithArguments("MethodWithArguments", 1))
}

func TestSpy_NumberOfCallWithArguments_SeveralCalls(t *testing.T) {
	tt := new(testing.T)
	spy := New[SpyExample](tt)

	spy.MethodWithArguments(1, "2", 3.0)
	spy.MethodWithArguments(1, "2", 3.0)
	spy.MethodWithArguments(1, "2", 3.0)

	assert.Equal(t, 3, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
}
