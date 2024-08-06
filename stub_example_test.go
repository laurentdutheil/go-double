package main

import (
	"fmt"
	"github.com/laurentdutheil/go-double/double"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
In production code
*/
type InterfaceExample interface {
	DoSomething(number int) (int, error)
}

// System Under Test
type SUT struct {
	dependency InterfaceExample
}

func (s SUT) MethodToTest(number int) error {
	n, err := s.dependency.DoSomething(number)
	if err != nil {
		return err
	}
	// do something with the result of dependency
	n = n % 100

	return nil
}

/*
In test file
*/
type MyStubObject struct {
	double.Stub
}

func (m *MyStubObject) DoSomething(number int) (int, error) {
	args := m.Called(number)
	return args.Int(0), args.Error(1)
}

func TestExample_Stub(t *testing.T) {
	t.Run("with On", func(t *testing.T) {
		stub := double.New[MyStubObject](t)
		stub.On("DoSomething", 3).Return(4, nil)
		// the next line is not called but does not fail the test
		stub.On("DoSomething", 123).Return(-1, fmt.Errorf("stub error"))

		sut := SUT{stub}
		err := sut.MethodToTest(3)

		assert.Nil(t, err)
	})

	t.Run("with When to have typing", func(t *testing.T) {
		stub := double.New[MyStubObject](t)
		stub.When(stub.DoSomething, 3).Return(4, nil)
		// the next line is not called but does not fail the test
		stub.When(stub.DoSomething, 123).Return(-1, fmt.Errorf("stub error"))

		sut := SUT{stub}
		err := sut.MethodToTest(3)

		assert.Nil(t, err)
	})
}
