package double_test

import (
	"fmt"
	. "github.com/laurentdutheil/go-double/double"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStub(t *testing.T) {
	tests := []struct {
		name        string
		constructor func(t TestingT) InterfaceTestStub
	}{
		{"for stub", func(t TestingT) InterfaceTestStub { return New[StubExample](t) }},
		{"for spy", func(t TestingT) InterfaceTestStub { return New[SpyExample](t) }},
		{"for mock", func(t TestingT) InterfaceTestStub { return New[MockExample](t) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			t.Run("On", func(t *testing.T) {
				t.Run("Predefine method name", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)

					call := stub.On("Method")

					assert.Equal(t, "Method", call.MethodName)
					assert.Contains(t, stub.PredefinedCalls(), call)
				})

				t.Run("Predefine method name and arguments", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)

					call := stub.On("MethodWithArguments", 1, "2", 3.0)

					assert.Equal(t, "MethodWithArguments", call.MethodName)
					assert.Contains(t, stub.PredefinedCalls(), call)
					assert.Len(t, call.Arguments, 3)
					assert.Contains(t, call.Arguments, 1)
					assert.Contains(t, call.Arguments, "2")
					assert.Contains(t, call.Arguments, 3.0)
				})
			})

			t.Run("Called", func(t *testing.T) {
				t.Run("Panic if do use the New constructor method incorrectly", func(t *testing.T) {
					stub := test.constructor(nil)

					expectedMessage := "Please use double.New constructor to initialize correctly."
					assert.PanicsWithValue(t, expectedMessage, func() { stub.Method() })
				})

				t.Run("FailNow when private method is used with Called", func(t *testing.T) {
					st := &SpiedTestingT{}
					stub := test.constructor(st)
					stub.On("privateMethod").Return(nil)

					st.AssertFailNowWasCalled(t, func() {
						_ = stub.privateMethod()
					})

					assert.Equal(t, "couldn't get the caller method information. 'privateMethod' is private or does not exist\n\tUse MethodCalled instead of Called in stub implementation.", st.errorMessages[0])
				})
			})

			t.Run("MethodCalled", func(t *testing.T) {
				t.Run("Panic if do use the New constructor method incorrectly", func(t *testing.T) {
					stub := test.constructor(nil)

					expectedMessage := "Please use double.New constructor to initialize correctly."
					assert.PanicsWithValue(t, expectedMessage, func() { _ = stub.privateMethodWithMethodCalled(1) })
				})

				t.Run("Use MethodCalled on private method", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)
					expectedErr := fmt.Errorf("stubbed error")
					stub.On("privateMethodWithMethodCalled", 1).Return(expectedErr)

					err := stub.privateMethodWithMethodCalled(1)

					assert.Equal(t, expectedErr, err)
				})
			})

			t.Run("On Return", func(t *testing.T) {
				t.Run("Predefine return arguments", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)
					expectedInt := 1
					expectedErr := fmt.Errorf("stubbed error")
					stub.On("MethodWithReturnArguments").Return(expectedInt, expectedErr)

					aInt, err := stub.MethodWithReturnArguments()

					assert.Equal(t, expectedInt, aInt)
					assert.Equal(t, expectedErr, err)
				})

				t.Run("Predefine return arguments with arguments checking", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)
					expectedInt := 1
					expectedErr := fmt.Errorf("stubbed error")
					stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)

					aInt, err := stub.MethodWithArgumentsAndReturnArguments(123, "123", 123.0)

					assert.Equal(t, expectedInt, aInt)
					assert.Equal(t, expectedErr, err)
				})

				t.Run("FailNow when arguments don't match", func(t *testing.T) {
					st := &SpiedTestingT{}
					stub := test.constructor(st)
					expectedInt := 1
					expectedErr := fmt.Errorf("stubbed error")
					stub.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(expectedInt, expectedErr)

					st.AssertFailNowWasCalled(t, func() {
						_, _ = stub.MethodWithArgumentsAndReturnArguments(12, "", 1.0)
					})
					assert.Equal(t, "I don't know what to return because the method call was unexpected.\n\tDo Stub.On(\"MethodWithArgumentsAndReturnArguments\").Return(...) first", st.errorMessages[0])
				})

				t.Run("Don't panic when method have no return arguments. Even if there is no predefined call", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)

					assert.NotPanics(t, func() { stub.Method() })
				})

			})

			t.Run("Times", func(t *testing.T) {
				t.Run("Return predefined return arguments once. And FailNow on the additional call", func(t *testing.T) {
					st := &SpiedTestingT{}
					stub := test.constructor(st)
					expectedInt := 1
					expectedErr := fmt.Errorf("stubbed error")
					stub.On("MethodWithReturnArguments").Return(expectedInt, expectedErr).Once()

					aInt, err := stub.MethodWithReturnArguments()
					assert.Equal(t, expectedInt, aInt)
					assert.Equal(t, expectedErr, err)

					st.AssertFailNowWasCalled(t, func() {
						_, _ = stub.MethodWithReturnArguments()
					})
				})

				t.Run("Return predefined return arguments twice. And panic on the additional call", func(t *testing.T) {
					st := &SpiedTestingT{}
					stub := test.constructor(st)
					expectedErr := fmt.Errorf("stubbed error")
					stub.On("MethodWithReturnArguments").Return(1, expectedErr).Twice()

					for i := 1; i <= 2; i++ {
						aInt, err := stub.MethodWithReturnArguments()
						assert.Equal(t, 1, aInt)
						assert.Equal(t, expectedErr, err)
					}

					st.AssertFailNowWasCalled(t, func() {
						_, _ = stub.MethodWithReturnArguments()
					})
				})

				t.Run("Return predefined return arguments n times. And panic on the additional call", func(t *testing.T) {
					st := &SpiedTestingT{}
					stub := test.constructor(st)
					expectedErr := fmt.Errorf("stubbed error")
					stub.On("MethodWithReturnArguments").Return(1, expectedErr).Times(4)

					for i := 1; i <= 4; i++ {
						aInt, err := stub.MethodWithReturnArguments()
						assert.Equal(t, 1, aInt)
						assert.Equal(t, expectedErr, err)
					}

					st.AssertFailNowWasCalled(t, func() {
						_, _ = stub.MethodWithReturnArguments()
					})
				})

				t.Run("Return different predefined return arguments. And panic on the additional call", func(t *testing.T) {
					st := &SpiedTestingT{}
					stub := test.constructor(st)
					expectedErr := fmt.Errorf("stubbed error")
					stub.On("MethodWithReturnArguments").Return(1, expectedErr).Once()
					stub.On("MethodWithReturnArguments").Return(2, expectedErr).Once()

					for i := 1; i <= 2; i++ {
						aInt, err := stub.MethodWithReturnArguments()
						assert.Equal(t, i, aInt)
						assert.Equal(t, expectedErr, err)
					}

					st.AssertFailNowWasCalled(t, func() {
						_, _ = stub.MethodWithReturnArguments()
					})
				})
			})

			t.Run("On Panic", func(t *testing.T) {
				t.Run("Panic when predefined call say so", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)
					stub.On("Method").Panic("panic message for example method")

					assert.PanicsWithValue(t, "panic message for example method", func() {
						stub.Method()
					})
				})
			})

			t.Run("On WailUntil", func(t *testing.T) {
				t.Run("Wait until the duration", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)
					stub.On("Method").WaitUntil(time.After(10 * time.Millisecond))

					done := make(chan string)
					go func() {
						stub.Method()
						done <- "done"
					}()

					// check that it is not done before
					select {
					case <-time.After(5 * time.Millisecond):
						// Pass
					case <-done:
						assert.Fail(t, "Have to wait until the duration")
					}

					// check that it is done after
					select {
					case <-time.After(20 * time.Millisecond):
						assert.Fail(t, "The wait is too long")
					case msg := <-done:
						assert.Equal(t, "done", msg)
					}
				})
			})

			t.Run("On After", func(t *testing.T) {
				t.Run("Wait until the duration", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)
					stub.On("Method").After(10 * time.Millisecond)

					done := make(chan string)
					go func() {
						stub.Method()
						done <- "done"
					}()

					// check that it is not done before
					select {
					case <-time.After(5 * time.Millisecond):
						// Pass
					case <-done:
						assert.Fail(t, "Have to wait until the duration")
					}

					// check that it is done after
					select {
					case <-time.After(20 * time.Millisecond):
						assert.Fail(t, "The wait is too long")
					case msg := <-done:
						assert.Equal(t, "done", msg)
					}
				})
			})

			t.Run("On Run", func(t *testing.T) {
				t.Run("Run function on a called method without argument", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)
					var called = false
					fn := func(_ Arguments) {
						called = true
					}
					stub.On("Method").Run(fn)

					stub.Method()

					assert.True(t, called)
				})

				t.Run("Run function on a called method with argument", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)
					fn := func(args Arguments) {
						arg := args[0].(*ExampleType)
						arg.ran = true
					}
					stub.On("MethodWithReferenceArgument", AnythingOfType("*double_test.ExampleType")).Run(fn)

					ref := &ExampleType{}
					stub.MethodWithReferenceArgument(ref)

					assert.True(t, ref.ran)
				})
			})

			t.Run("When", func(t *testing.T) {
				t.Run("Predefine method name", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)

					call := stub.When(stub.Method)

					assert.Equal(t, "Method", call.MethodName)
					assert.Contains(t, stub.PredefinedCalls(), call)
				})

				t.Run("Predefine method name and arguments", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)

					call := stub.When(stub.MethodWithArguments, 1, "2", 3.0)

					assert.Equal(t, "MethodWithArguments", call.MethodName)
					assert.Contains(t, stub.PredefinedCalls(), call)
					assert.Len(t, call.Arguments, 3)
					assert.Contains(t, call.Arguments, 1)
					assert.Contains(t, call.Arguments, "2")
					assert.Contains(t, call.Arguments, 3.0)
				})

				t.Run("Panic if pass anything other than a function ", func(t *testing.T) {
					tt := new(testing.T)
					stub := test.constructor(tt)

					expectedMessage := "Please pass the function as an argument : stub.When(stub.Method)"
					assert.PanicsWithValue(t, expectedMessage, func() { stub.When("not a function", 1, "2", 3.0) })
				})
			})
		})
	}
}
