package double_test

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

type TestStub struct {
	Stub
}

func (s *TestStub) Method() {
	s.Called()
}

func (s *TestStub) MethodWithArgs(aInt int, aString string, aFloat float64) {
	s.Called(aInt, aString, aFloat)
}

func (s *TestStub) MethodWithReturnArgs() (int, error) {
	return 0, nil
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

func TestStubOn_RegisterMethodNameWithReturnArgs(t *testing.T) {
	stub := &TestStub{}

	call := stub.On("Method").Return(1, nil)

	assert.Equal(t, "Method", call.MethodName)
	assert.Contains(t, stub.RegisteredCalls, call)
	assert.Contains(t, call.ReturnArguments, 1)
	assert.Contains(t, call.ReturnArguments, nil)
}

type TestInterface interface {
	Method()
	MethodWithArgs(aInt int, aString string, aFloat float64)
	MethodWithReturnArgs() (int, error)
}

type TestSUT struct {
	dependency TestInterface
}

func (sut TestSUT) method() {
	sut.dependency.Method()
}

func (sut TestSUT) methodWithArgs(aInt int) {
	sut.dependency.MethodWithArgs(aInt, strconv.Itoa(aInt), float64(aInt))
}

func TestStub_CallIsRecorded(t *testing.T) {
	stub := &TestStub{}
	sut := &TestSUT{stub}

	sut.method()

	assert.Len(t, stub.Calls, 1)
	assert.Equal(t, *NewCall("Method"), stub.Calls[0])
}

func TestStub_CallWithArgumentsIsRecorded(t *testing.T) {
	stub := &TestStub{}
	sut := &TestSUT{stub}

	sut.methodWithArgs(123)

	assert.Len(t, stub.Calls, 1)
	assert.Equal(t, *NewCall("MethodWithArgs", 123, "123", 123.0), stub.Calls[0])
}
