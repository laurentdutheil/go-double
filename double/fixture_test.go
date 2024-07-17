package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

type StubExample struct {
	Stub[StubExample]
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

func (s *StubExample) MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error) {
	arguments := s.Called(aInt, aString, aFloat)
	return arguments[0].(int), arguments[1].(error)
}

func (s *StubExample) MethodWithReferenceArgument(ref *ExampleType) {
	s.Called(ref)
}

type SpyExample struct {
	Spy[SpyExample]
}

func (s *SpyExample) Method() {
	s.Called()
}

func (s *SpyExample) MethodWithArguments(aInt int, aString string, aFloat float64) {
	s.Called(aInt, aString, aFloat)
}

func (s *SpyExample) MethodWithReturnArguments() (int, error) {
	arguments := s.Called()
	return arguments[0].(int), arguments[1].(error)
}

func (s *SpyExample) MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error) {
	arguments := s.Called(aInt, aString, aFloat)
	return arguments[0].(int), arguments[1].(error)
}

func (s *SpyExample) MethodWithReferenceArgument(ref *ExampleType) {
	s.Called(ref)
}

type MockExample struct {
	Mock[MockExample]
}

func (s *MockExample) Method() {
	s.Called()
}

func (s *MockExample) MethodWithOneArgument(aInt int) {
	s.Called(aInt)
}

func (s *MockExample) MethodWithArguments(aInt int, aString string, aFloat float64) {
	s.Called(aInt, aString, aFloat)
}

func (s *MockExample) MethodWithReturnArguments() (int, error) {
	arguments := s.Called()
	return arguments[0].(int), arguments[1].(error)
}

func (s *MockExample) MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error) {
	arguments := s.Called(aInt, aString, aFloat)
	return arguments[0].(int), arguments[1].(error)
}

func (s *MockExample) MethodWithReferenceArgument(ref *ExampleType) {
	s.Called(ref)
}

type ExampleType struct {
	ran bool
}

type SpiedTestingT struct {
	errorfFormat  string
	errorfArgs    []interface{}
	helperCalled  bool
	failNowCalled bool
}

func (s *SpiedTestingT) Errorf(format string, args ...interface{}) {
	s.errorfFormat = format
	s.errorfArgs = args
}

func (s *SpiedTestingT) Helper() {
	s.helperCalled = true
}

func (s *SpiedTestingT) FailNow() {
	s.failNowCalled = true
	panic("SpiedTestingT.FailNow() called")
}

func (s *SpiedTestingT) AssertFailNowWasCalled(t *testing.T, f func()) {
	assert.PanicsWithValue(t, "SpiedTestingT.FailNow() called", f)
	assert.True(t, s.failNowCalled)
}

// Check if SpiedTestingT implements all methods of TestingT
var _ TestingT = (*SpiedTestingT)(nil)
