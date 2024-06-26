package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

type TestInterface interface {
	Method()
	MethodWithArgs(aInt int, aString string, aFloat float64)
}

type TestStub struct {
	Stub
}

func (s *TestStub) MethodWithArgs(aInt int, aString string, aFloat float64) {
}

func (s *TestStub) Method() {
}

func TestStubOn_RegisterMethodName(t *testing.T) {
	stub := &TestStub{}

	call := stub.On("Method")

	assert.Equal(t, "Method", call.MethodName)
	assert.Contains(t, stub.RegisteredCalls, call)
}

func TestStubOn_RegisterMethodNameAndArgs(t *testing.T) {
	stub := &TestStub{}

	call := stub.On("MethodWithArgs", 1, "2", 3.0)

	assert.Equal(t, "MethodWithArgs", call.MethodName)
	assert.Contains(t, stub.RegisteredCalls, call)
	assert.Len(t, call.Arguments, 3)
	assert.Contains(t, call.Arguments, 1)
	assert.Contains(t, call.Arguments, "2")
	assert.Contains(t, call.Arguments, 3.0)
}
