package double_test

import (
	"github.com/laurentdutheil/go-double/double"
	"strconv"
)

type InterfaceExample interface {
	Method()
	MethodWithArguments(aInt int, aString string, aFloat float64)
	MethodWithReturnArguments() (int, error)
}

type StubExample struct {
	double.Stub
}

func (s *StubExample) Method() {
	s.Called()
}

func (s *StubExample) MethodWithArguments(aInt int, aString string, aFloat float64) {
	s.Called(aInt, aString, aFloat)
}

func (s *StubExample) MethodWithReturnArguments() (int, error) {
	arguments := s.Called()
	return arguments[0].(int), arguments[1].(error)
}

type SUTExample struct {
	dependency InterfaceExample
}

func (sut SUTExample) method() {
	sut.dependency.Method()
}

func (sut SUTExample) methodWithArguments(aInt int) {
	sut.dependency.MethodWithArguments(aInt, strconv.Itoa(aInt), float64(aInt))
}

func (sut SUTExample) methodWithReturnArguments() (int, error) {
	return sut.dependency.MethodWithReturnArguments()
}
