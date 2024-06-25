package double_test

import (
	"slices"
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

	if "Method" != call.MethodName {
		t.Errorf("Stub.On() should register Method name\n got: %s\nwant: %s", call.MethodName, "Method")
	}
	if !slices.Contains(stub.RegisteredCalls, call) {
		t.Errorf("Stub.RegistredCalls should contain the result of Stub.On()")
	}
}
