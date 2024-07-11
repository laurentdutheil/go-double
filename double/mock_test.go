package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestMock_AssertCall(t *testing.T) {

	t.Run("t.Helper is called", func(t *testing.T) {
		mock := MockExample{}
		st := &SpiedTestingT{}

		mock.AssertCalled(st, "Method")

		assert.True(t, st.helperCalled)
	})

	t.Run("Return false when method is not called", func(t *testing.T) {
		mock := MockExample{}
		tt := new(testing.T)

		result := mock.AssertCalled(tt, "Method")

		assert.False(t, result)
	})

	t.Run("Return true when method is called", func(t *testing.T) {
		mock := MockExample{}
		tt := new(testing.T)
		mock.Method()

		result := mock.AssertCalled(tt, "Method")

		assert.True(t, result)
	})

	t.Run("t.Errorf is called with right message on error", func(t *testing.T) {
		mock := MockExample{}
		st := &SpiedTestingT{}

		mock.AssertCalled(st, "Method")

		assert.Equal(t, "\n%s", st.errorfFormat)
		errorMessage := st.errorfArgs[0]
		assert.Contains(t, errorMessage, "Error Trace:")
		assert.Contains(t, errorMessage, "Should have called")
		assert.Contains(t, errorMessage, "Expected \"Method\" to have been called")
		assert.Contains(t, errorMessage, "but no actual calls happened")
	})

}

type SpiedTestingT struct {
	errorfFormat string
	errorfArgs   []interface{}
	helperCalled bool
}

func (s *SpiedTestingT) Errorf(format string, args ...interface{}) {
	s.errorfFormat = format
	s.errorfArgs = args
}

func (s *SpiedTestingT) Helper() {
	s.helperCalled = true
}

// Check if SpiedTestingT implements all methods of TestingT
var _ TestingT = (*SpiedTestingT)(nil)
