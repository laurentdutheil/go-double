package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestSpy(t *testing.T) {
	tests := []struct {
		name        string
		constructor func(t TestingT) InterfaceTestSpy
	}{
		{"for spy", func(t TestingT) InterfaceTestSpy { return New[SpyExample](t) }},
		{"for mock", func(t TestingT) InterfaceTestSpy { return New[MockExample](t) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Run("Called", func(t *testing.T) {
				t.Run("Register actual call", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.Method()

					actualCalls := spy.ActualCalls()
					assert.Len(t, actualCalls, 1)
					assert.Equal(t, NewActualCall("Method"), actualCalls[0])
				})

				t.Run("Register actual call with arguments", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.MethodWithArguments(123, "123", 123.0)

					actualCalls := spy.ActualCalls()
					assert.Len(t, actualCalls, 1)
					assert.Equal(t, NewActualCall("MethodWithArguments", 123, "123", 123.0), actualCalls[0])
				})
			})

			t.Run("MethodCalled", func(t *testing.T) {
				t.Run("Register actual call on private method", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)
					expectedErr := fmt.Errorf("stubbed error")
					spy.On("privateMethodWithMethodCalled", 1).Return(expectedErr)

					err := spy.privateMethodWithMethodCalled(1)

					assert.Equal(t, expectedErr, err)
					actualCalls := spy.ActualCalls()
					assert.Len(t, actualCalls, 1)
					assert.Equal(t, NewActualCall("privateMethodWithMethodCalled", 1), actualCalls[0])
				})
			})

			t.Run("AddActualCall", func(t *testing.T) {
				t.Run("Register actual call", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					result := spy.methodOnlyAddActualCall(1)

					assert.Equal(t, 123, result)
					actualCalls := spy.ActualCalls()
					assert.Len(t, actualCalls, 1)
					assert.Equal(t, NewActualCall("methodOnlyAddActualCall", 1), actualCalls[0])
				})
			})

			t.Run("NumberOfCalls", func(t *testing.T) {
				t.Run("Zero call", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					numberOfCalls := spy.NumberOfCalls("Method")

					assert.Equal(t, 0, numberOfCalls)
				})

				t.Run("Several calls", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.Method()
					spy.Method()

					numberOfCalls := spy.NumberOfCalls("Method")

					assert.Equal(t, 2, numberOfCalls)
				})

				t.Run("Several calls with different arguments", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.MethodWithArguments(1, "1", 1.2)
					spy.MethodWithArguments(2, "2", 2.2)

					numberOfCalls := spy.NumberOfCalls("MethodWithArguments")

					assert.Equal(t, 2, numberOfCalls)
				})
			})

			t.Run("NumberOfActualCallsWithArguments", func(t *testing.T) {
				t.Run("Zero call", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					numberOfCalls := spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0)

					assert.Equal(t, 0, numberOfCalls)
				})

				t.Run("Several calls", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.MethodWithArguments(1, "2", 3.0)
					spy.MethodWithArguments(1, "2", 3.0)
					spy.MethodWithArguments(2, "3", 4.5)

					numberOfCalls := spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0)

					assert.Equal(t, 2, numberOfCalls)
				})

				t.Run("One call with wrong arguments", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.MethodWithArguments(0, "2", 3.0)

					numberOfCalls := spy.NumberOfCallsWithArguments("MethodWithArguments", 1, "2", 3.0)

					assert.Equal(t, 0, numberOfCalls)
				})

				t.Run("One call with wrong number of arguments", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.MethodWithArguments(0, "2", 3.0)

					numberOfCalls := spy.NumberOfCallsWithArguments("MethodWithArguments", 1)

					assert.Equal(t, 0, numberOfCalls)
				})

				t.Run("Several calls with Anything", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.MethodWithArguments(1, "2", 3.0)
					spy.MethodWithArguments(2, "2", 3.0)
					spy.MethodWithArguments(3, "2", 3.0)

					numberOfCalls := spy.NumberOfCallsWithArguments("MethodWithArguments", Anything, "2", 3.0)

					assert.Equal(t, 3, numberOfCalls)
				})

				t.Run("Several calls with AnythingOfType", func(t *testing.T) {
					tt := new(testing.T)
					spy := test.constructor(tt)

					spy.MethodWithArguments(1, "2", 3.0)
					spy.MethodWithArguments(2, "2", 3.0)
					spy.MethodWithArguments(3, "2", 3.0)

					numberOfCalls := spy.NumberOfCallsWithArguments("MethodWithArguments", AnythingOfType("int"), "2", 3.0)

					assert.Equal(t, 3, numberOfCalls)
				})
			})
		})
	}
}
