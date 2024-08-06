package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

type StubExample struct {
	Stub
}

func (s *StubExample) Method() {
	s.Called()
}

func (s *StubExample) MethodWithArguments(aInt int, aString string, aFloat float64) {
	s.Called(aInt, aString, aFloat)
}

func (s *StubExample) MethodWithReturnArguments() (int, error) {
	arguments := s.Called()
	return arguments.Int(0), arguments.Error(1)
}

func (s *StubExample) MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error) {
	arguments := s.Called(aInt, aString, aFloat)
	return arguments.Int(0), arguments.Error(1)
}

func (s *StubExample) MethodWithReferenceArgument(ref *ExampleType) {
	s.Called(ref)
}

func (s *StubExample) privateMethod() error {
	arguments := s.Called()
	return arguments.Error(0)
}

func (s *StubExample) privateMethodWithMethodCalled(aInt int) error {
	methodInformation := MethodInformation{
		Name:   "privateMethodWithMethodCalled",
		NumOut: 1,
	}
	arguments := s.MethodCalled(methodInformation, aInt)
	return arguments.Error(0)
}

type SpyExample struct {
	Spy
}

func (s *SpyExample) Method() {
	s.Called()
}

func (s *SpyExample) MethodWithArguments(aInt int, aString string, aFloat float64) {
	s.Called(aInt, aString, aFloat)
}

func (s *SpyExample) MethodWithReturnArguments() (int, error) {
	arguments := s.Called()
	return arguments.Int(0), arguments.Error(1)
}

func (s *SpyExample) MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error) {
	arguments := s.Called(aInt, aString, aFloat)
	return arguments.Int(0), arguments.Error(1)
}

func (s *SpyExample) MethodWithReferenceArgument(ref *ExampleType) {
	s.Called(ref)
}

func (s *SpyExample) privateMethod() error {
	arguments := s.Called()
	return arguments.Error(0)
}

func (s *SpyExample) privateMethodWithMethodCalled(aInt int) error {
	methodInformation := MethodInformation{
		Name:   "privateMethodWithMethodCalled",
		NumOut: 1,
	}
	arguments := s.MethodCalled(methodInformation, aInt)
	return arguments.Error(0)
}

func (s *SpyExample) methodOnlyAddActualCall(aInt int) int {
	s.AddActualCall(aInt)
	return 123
}

type MockExample struct {
	Mock
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
	return arguments.Int(0), arguments.Error(1)
}

func (s *MockExample) MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error) {
	arguments := s.Called(aInt, aString, aFloat)
	return arguments.Int(0), arguments.Error(1)
}

func (s *MockExample) MethodWithReferenceArgument(ref *ExampleType) {
	s.Called(ref)
}

func (s *MockExample) privateMethod() error {
	arguments := s.Called()
	return arguments.Error(0)
}

func (s *MockExample) privateMethodWithMethodCalled(aInt int) error {
	methodInformation := MethodInformation{
		Name:   "privateMethodWithMethodCalled",
		NumOut: 1,
	}
	arguments := s.MethodCalled(methodInformation, aInt)
	return arguments.Error(0)
}

func (s *MockExample) methodOnlyAddActualCall(aInt int) int {
	s.AddActualCall(aInt)
	return 123
}

type ExampleType struct {
	ran bool
}

type SpiedTestingT struct {
	errorMessages []string
	helperCalled  bool
	failNowCalled bool
}

func (s *SpiedTestingT) Errorf(format string, args ...interface{}) {
	s.errorMessages = append(s.errorMessages, fmt.Sprintf(format, args...))
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

type InterfaceStubExample interface {
	Method()
	MethodWithArguments(aInt int, aString string, aFloat float64)
	MethodWithReturnArguments() (int, error)
	MethodWithArgumentsAndReturnArguments(aInt int, aString string, aFloat float64) (int, error)
	MethodWithReferenceArgument(ref *ExampleType)
	privateMethod() error
	privateMethodWithMethodCalled(aInt int) error
}

// Check if StubExample implements all methods of InterfaceStubExample
var _ InterfaceStubExample = (*StubExample)(nil)

// Check if SpyExample implements all methods of InterfaceStubExample
var _ InterfaceStubExample = (*SpyExample)(nil)

// Check if MockExample implements all methods of InterfaceStubExample
var _ InterfaceStubExample = (*MockExample)(nil)

type InterfaceTestStub interface {
	InterfaceStubExample
	IStub
}

type InterfaceSpyExample interface {
	InterfaceStubExample
	methodOnlyAddActualCall(aInt int) int
}

type InterfaceTestSpy interface {
	InterfaceSpyExample
	ISpy
}
