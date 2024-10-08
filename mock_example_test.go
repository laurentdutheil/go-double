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
type Interface interface {
	GetSomething(number int) (int, error)
	DoSomething(number int)
}

type ObjectToTest struct {
	dependency Interface
}

func (o ObjectToTest) MethodToTest(number int) error {
	n, err := o.dependency.GetSomething(number)
	if err != nil {
		return err
	}
	n = (n + 42) % 10
	o.dependency.DoSomething(n)
	return nil
}

/*
In test file
*/
type MyMockObject struct {
	double.Mock
}

func (m *MyMockObject) GetSomething(number int) (int, error) {
	args := m.Called(number)
	return args.Int(0), args.Error(1)
}

func (m *MyMockObject) DoSomething(number int) {
	m.Called(number)
}

func TestExample_Mock(t *testing.T) {
	t.Run("Old fashion way with AssertExpectations", func(t *testing.T) {
		mock := double.New[MyMockObject](t)
		mock.On("GetSomething", 3).Return(4, nil)
		mock.On("DoSomething", 6)

		objectToTest := ObjectToTest{mock}
		_ = objectToTest.MethodToTest(3)

		mock.AssertExpectations(t)
	})

	t.Run("Stub the requests and mock the command", func(t *testing.T) {
		mock := double.New[MyMockObject](t)
		mock.On("GetSomething", 3).Return(4, nil)

		objectToTest := ObjectToTest{mock}
		_ = objectToTest.MethodToTest(3)

		mock.AssertCalled(t, "DoSomething", 6)
	})

	t.Run("with AssertNotCalled", func(t *testing.T) {
		mock := double.New[MyMockObject](t)
		expectedError := fmt.Errorf("mock error")
		mock.On("GetSomething", double.Anything).Return(0, expectedError)

		objectToTest := ObjectToTest{mock}
		err := objectToTest.MethodToTest(3)

		assert.Equal(t, expectedError, err)
		mock.AssertNotCalled(t, "DoSomething", double.Anything)
	})

	t.Run("Assert the order of the calls", func(t *testing.T) {
		mock1 := double.New[MyMockObject](t)
		mock1.On("GetSomething", double.Anything).Return(1, nil)
		mock2 := double.New[MyMockObject](t)
		mock2.On("GetSomething", double.Anything).Return(2, nil)
		inOrder := double.InOrder(mock1, mock2)
		objectToTest1 := ObjectToTest{mock1}
		objectToTest2 := ObjectToTest{mock2}

		_ = objectToTest1.MethodToTest(1)
		_ = objectToTest2.MethodToTest(2)

		inOrder.AssertCalled(t, mock1, "GetSomething", 1)
		inOrder.AssertCalled(t, mock1, "DoSomething", double.AnythingOfType("int"))
		inOrder.AssertCalled(t, mock2, "GetSomething", 2)
		inOrder.AssertCalled(t, mock2, "DoSomething", double.AnythingOfType("int"))
		inOrder.AssertNoMoreExpectations(t)
	})
}
