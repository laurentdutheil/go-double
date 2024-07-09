package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestSpy_Called_RegisterActualCall(t *testing.T) {
	spy := &SpyExample{}
	sut := &SUTExample{spy}

	sut.method()

	assert.Len(t, spy.ActualCalls, 1)
	assert.Equal(t, *NewCall("Method"), spy.ActualCalls[0])
}

func TestSpy_Called_RegisterActualCallWithArguments(t *testing.T) {
	spy := &SpyExample{}
	sut := &SUTExample{spy}

	sut.methodWithArguments(123)

	assert.Len(t, spy.ActualCalls, 1)
	assert.Equal(t, *NewCall("MethodWithArguments", 123, "123", 123.0), spy.ActualCalls[0])
}

func TestSpy_NumberOfCall_ZeroCall(t *testing.T) {
	spy := &SpyExample{}

	assert.Equal(t, 0, spy.NumberOfCall("Method"))
}

func TestSpy_NumberOfCall_SeveralCalls(t *testing.T) {
	spy := &SpyExample{}

	spy.Method()
	spy.Method()

	assert.Equal(t, 2, spy.NumberOfCall("Method"))
}