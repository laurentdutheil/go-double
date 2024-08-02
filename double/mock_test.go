package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestMock(t *testing.T) {
	t.Run("AssertNumberOfCalls", func(t *testing.T) {
		t.Run("t.Helper is called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)

			mock.AssertNumberOfCalls(st, "Method", 1)

			assert.True(t, st.helperCalled)
		})

		t.Run("Return false when number of calls is incorrect", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)

			result := mock.AssertNumberOfCalls(tt, "Method", 1)

			assert.False(t, result)
		})

		t.Run("Return true when number of calls is correct", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)
			mock.Method()

			result := mock.AssertNumberOfCalls(tt, "Method", 1)

			assert.True(t, result)
		})

		t.Run("t.Errorf is called with right message when number of calls is incorrect", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)
			mock.Method()

			mock.AssertNumberOfCalls(st, "Method", 2)

			assert.Len(t, st.errorMessages, 1)
			assert.Contains(t, st.errorMessages[0], "Expected number of calls (2) does not match the actual number of calls (1).")
		})
	})

	t.Run("AssertNumberOfCallsWithArguments", func(t *testing.T) {
		t.Run("t.Helper is called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)

			mock.AssertNumberOfCallsWithArguments(st, 1, "MethodWithArguments", 1, "2", 3.4)

			assert.True(t, st.helperCalled)
		})

		t.Run("Return false when method is not called", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)

			result := mock.AssertNumberOfCallsWithArguments(tt, 1, "MethodWithArguments", 1, "2", 3.4)

			assert.False(t, result)
		})

		t.Run("Return true when method is not called", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)

			mock.MethodWithArguments(1, "2", 3.4)

			result := mock.AssertNumberOfCallsWithArguments(tt, 1, "MethodWithArguments", 1, "2", 3.4)

			assert.True(t, result)
		})

		t.Run("t.Errorf is called with right message when number of calls is incorrect", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)
			mock.MethodWithArguments(1, "2", 3.4)

			mock.AssertNumberOfCallsWithArguments(st, 2, "MethodWithArguments", 1, "2", 3.4)

			assert.Len(t, st.errorMessages, 1)
			assert.Contains(t, st.errorMessages[0], "Expected number of calls (2) does not match the actual number of calls (1).")
		})
	})

	t.Run("AssertCall", func(t *testing.T) {
		t.Run("t.Helper is called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)

			mock.AssertCalled(st, "Method")

			assert.True(t, st.helperCalled)
		})

		t.Run("Return false when method is not called", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)

			result := mock.AssertCalled(tt, "Method")

			assert.False(t, result)
		})

		t.Run("Return true when method is called", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)
			mock.Method()

			result := mock.AssertCalled(tt, "Method")

			assert.True(t, result)
		})

		t.Run("t.Errorf is called with right message when method is not called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)

			mock.AssertCalled(st, "Method")

			assert.Len(t, st.errorMessages, 1)
			assert.Contains(t, st.errorMessages[0], "Should have called with given arguments\n\tMessages:   \tExpected \"Method\" to have been called with:\n\t            \t[]\n\t            \tbut no actual calls happened\n")
		})

		t.Run("Return true when method is called with right arguments", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)
			mock.MethodWithArguments(1, "1", 1.0)

			result := mock.AssertCalled(tt, "MethodWithArguments", 1, "1", 1.0)

			assert.True(t, result)
		})

		t.Run("Return false when method is called with wrong arguments", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)
			mock.MethodWithArguments(2, "1", 1.0)

			result := mock.AssertCalled(tt, "MethodWithArguments", 1, "1", 1.0)

			assert.False(t, result)
		})

		t.Run("t.Errorf is called with right message when method is called with wrong arguments", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)
			mock.MethodWithArguments(2, "1", 1.0)
			mock.MethodWithArguments(1, "3", 1.2)
			mock.MethodWithOneArgument(4)

			mock.AssertCalled(st, "MethodWithArguments", 1, "1", 1.0)

			assert.Len(t, st.errorMessages, 1)
			assert.Contains(t, st.errorMessages[0], "Should have called with given arguments\n\tMessages:   \tExpected \"MethodWithArguments\" to have been called with:\n\t            \t[1 1 1]\n\t            \tbut actual calls were:\n\t            \t        [2 1 1]\n\t            \t[1 3 1.2]\n")
		})
	})

	t.Run("AssertNotCalled", func(t *testing.T) {
		t.Run("t.Helper is called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)

			mock.AssertNotCalled(st, "Method")

			assert.True(t, st.helperCalled)
		})

		t.Run("Return false when method is called", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)
			mock.Method()

			result := mock.AssertNotCalled(tt, "Method")

			assert.False(t, result)
		})

		t.Run("Return true when method is not called", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)

			result := mock.AssertNotCalled(tt, "Method")

			assert.True(t, result)
		})

		t.Run("t.Errorf is called with right message when method is called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)
			mock.Method()

			mock.AssertNotCalled(st, "Method")

			assert.Len(t, st.errorMessages, 1)
			assert.Contains(t, st.errorMessages[0], "Should not have called with given arguments\n\tMessages:   \tExpected \"Method\" to not have been called with:\n\t            \t[]\n\t            \tbut actually it was.\n")
		})

		t.Run("Return true when method is called with other arguments", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)
			mock.MethodWithArguments(2, "1", 1.0)

			result := mock.AssertNotCalled(tt, "MethodWithArguments", 1, "1", 1.0)

			assert.True(t, result)
		})

		t.Run("Return false when method is called with same arguments", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)
			mock.MethodWithArguments(1, "1", 1.0)

			result := mock.AssertNotCalled(tt, "MethodWithArguments", 1, "1", 1.0)

			assert.False(t, result)
		})

		t.Run("t.Errorf is called with right message when method is called with same arguments", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)
			mock.MethodWithArguments(1, "3", 1.2)

			mock.AssertNotCalled(st, "MethodWithArguments", 1, "3", 1.2)

			assert.Len(t, st.errorMessages, 1)
			assert.Contains(t, st.errorMessages[0], "Should not have called with given arguments\n\tMessages:   \tExpected \"MethodWithArguments\" to not have been called with:\n\t            \t[1 3 1.2]\n\t            \tbut actually it was.\n")
		})
	})

	t.Run("AssertExpectations", func(t *testing.T) {
		t.Run("t.Helper is called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)

			mock.AssertExpectations(st)

			assert.True(t, st.helperCalled)
		})

		t.Run("Return true when there is no expectation", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)

			result := mock.AssertExpectations(tt)

			assert.True(t, result)
		})

		t.Run("Return false when one expectation is not called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)
			mock.On("MethodWithReturnArguments").Return(1, fmt.Errorf("expected Error"))

			result := mock.AssertExpectations(st)

			assert.False(t, result)
			assert.Len(t, st.errorMessages, 1)
			assert.Contains(t, st.errorMessages[0], "Should have called with given arguments\n\tMessages:   \tExpected \"MethodWithReturnArguments\" to have been called with:\n\t            \t[]\n\t            \tbut no actual calls happened\n")
		})

		t.Run("Return true when one expectation is called", func(t *testing.T) {
			tt := new(testing.T)
			mock := New[MockExample](tt)
			mock.On("MethodWithReturnArguments").Return(1, fmt.Errorf("expected Error"))

			_, _ = mock.MethodWithReturnArguments()
			result := mock.AssertExpectations(tt)

			assert.True(t, result)
		})

		t.Run("Return false when several expectations are not called", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)
			mock.On("MethodWithReturnArguments").Return(1, fmt.Errorf("expected Error"))
			mock.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(1, fmt.Errorf("expected Error"))

			result := mock.AssertExpectations(st)

			assert.False(t, result)
			assert.Len(t, st.errorMessages, 2)
			assert.Contains(t, st.errorMessages[0], "Should have called with given arguments\n\tMessages:   \tExpected \"MethodWithReturnArguments\" to have been called with:\n\t            \t[]\n\t            \tbut no actual calls happened\n")
			assert.Contains(t, st.errorMessages[1], "Should have called with given arguments\n\tMessages:   \tExpected \"MethodWithArgumentsAndReturnArguments\" to have been called with:\n\t            \t[123 123 123]\n\t            \tbut no actual calls happened\n")
		})

		t.Run("Return false when expectation are not called enough times", func(t *testing.T) {
			st := &SpiedTestingT{}
			mock := New[MockExample](st)
			mock.On("MethodWithArgumentsAndReturnArguments", 123, "123", 123.0).Return(1, fmt.Errorf("expected Error")).Times(2)

			_, _ = mock.MethodWithArgumentsAndReturnArguments(123, "123", 123.0)
			result := mock.AssertExpectations(st)

			assert.False(t, result)
			assert.Len(t, st.errorMessages, 1)
			assert.Contains(t, st.errorMessages[0], "Should have called with given arguments\n\tMessages:   \tExpected \"MethodWithArgumentsAndReturnArguments\" to have been called 2 times with:\n\t            \t[123 123 123]\n\t            \tbut actually it was called 1 times.\n")
		})
	})

	t.Run("Race condition", func(t *testing.T) {
		t.Run("Concurrent Return and Called", func(t *testing.T) {
			iterations := 1000
			tt := new(testing.T)
			m := New[MockExample](tt)

			call := m.On("Method")

			wg := sync.WaitGroup{}
			wg.Add(2)

			go func() {
				for i := 0; i < iterations; i++ {
					call.Return(10)
				}
				wg.Done()
			}()
			go func() {
				for i := 0; i < iterations; i++ {
					m.Method()
				}
				wg.Done()
			}()
			wg.Wait()
		})
	})
}
