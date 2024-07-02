package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestSpy_Called_RegisterActualCall(t *testing.T) {
	spy := &SpyExample{}
	spy.On("Method")
	sut := &SUTExample{spy}

	sut.method()

	assert.Len(t, spy.ActualCalls, 1)
	assert.Equal(t, *NewCall("Method"), spy.ActualCalls[0])
}

func TestSpy_Called_RegisterActualCallWithArguments(t *testing.T) {
	spy := &SpyExample{}
	spy.On("MethodWithArguments", 123, "123", 123.0)
	sut := &SUTExample{spy}

	sut.methodWithArguments(123)

	assert.Len(t, spy.ActualCalls, 1)
	assert.Equal(t, *NewCall("MethodWithArguments", 123, "123", 123.0), spy.ActualCalls[0])
}
