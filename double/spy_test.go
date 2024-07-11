package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestSpy_Called_RegisterActualCall(t *testing.T) {
	spy := &SpyExample{}

	spy.Method()

	assert.Len(t, spy.ActualCalls, 1)
	assert.Equal(t, *NewCall("Method"), spy.ActualCalls[0])
}

func TestSpy_Called_RegisterActualCallWithArguments(t *testing.T) {
	spy := &SpyExample{}

	spy.MethodWithArguments(123, "123", 123.0)

	assert.Len(t, spy.ActualCalls, 1)
	assert.Equal(t, *NewCall("MethodWithArguments", 123, "123", 123.0), spy.ActualCalls[0])
}

func TestSpy_NumberOfCall_ZeroCall(t *testing.T) {
	spy := &SpyExample{}

	assert.Equal(t, 0, spy.NumberOfCalls("Method"))
}

func TestSpy_NumberOfCall_SeveralCalls(t *testing.T) {
	spy := &SpyExample{}

	spy.Method()
	spy.Method()

	assert.Equal(t, 2, spy.NumberOfCalls("Method"))
}

func TestSpy_NumberOfCallWithArguments_ZeroCall(t *testing.T) {
	spy := &SpyExample{}

	assert.Equal(t, 0, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
}

func TestSpy_NumberOfCallWithArguments_OneCallWithWrongArguments(t *testing.T) {
	spy := &SpyExample{}

	spy.MethodWithArguments(0, "2", 3.0)

	assert.Equal(t, 0, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
}

func TestSpy_NumberOfCallWithArguments_SeveralCalls(t *testing.T) {
	spy := &SpyExample{}

	spy.MethodWithArguments(1, "2", 3.0)
	spy.MethodWithArguments(1, "2", 3.0)
	spy.MethodWithArguments(1, "2", 3.0)

	assert.Equal(t, 3, spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0))
}
