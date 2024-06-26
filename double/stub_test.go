package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

type TestInterface interface {
	Method()
}

type TestStub struct {
	Stub
}

func (s *TestStub) Method() {
}

func TestStubOn_RegisterMethodName(t *testing.T) {
	stub := &TestStub{}

	call := stub.On("Method")

	assert.Equal(t, "Method", call.MethodName)
	assert.Contains(t, stub.RegisteredCalls, call)
}
