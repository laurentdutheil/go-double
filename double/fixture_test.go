package double_test

import (
	"github.com/laurentdutheil/go-double/double"
)

type StubExample struct {
	double.Stub
}

func (s *StubExample) Method() {
	s.Called(s)
}

func (s *StubExample) MethodWithArguments(aInt int, aString string, aFloat float64) {
	s.Called(s, aInt, aString, aFloat)
}

func (s *StubExample) MethodWithReturnArguments() (int, error) {
	arguments := s.Called(s)
	return arguments[0].(int), arguments[1].(error)
}

func (s *StubExample) MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error) {
	arguments := s.Called(s, aInt, aString, aFloat)
	return arguments[0].(int), arguments[1].(error)
}

type SpyExample struct {
	double.Spy
}

func (s *SpyExample) Method() {
	s.Called(s)
}

func (s *SpyExample) MethodWithArguments(aInt int, aString string, aFloat float64) {
	s.Called(s, aInt, aString, aFloat)
}

func (s *SpyExample) MethodWithReturnArguments() (int, error) {
	arguments := s.Called(s)
	return arguments[0].(int), arguments[1].(error)
}

func (s *SpyExample) MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error) {
	arguments := s.Called(s, aInt, aString, aFloat)
	return arguments[0].(int), arguments[1].(error)
}
