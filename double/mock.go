package double

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Mock struct {
	Spy
}

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
	Helper()
}

// Check if TestingT interface can wrap testing.T
var _ TestingT = (*testing.T)(nil)

func (m *Mock) AssertCalled(t TestingT, methodName string) bool {
	t.Helper()

	result := m.Spy.NumberOfCall(methodName) > 0
	if !result {
		assert.Fail(t, "Should have called",
			fmt.Sprintf("Expected %q to have been called\nbut no actual calls happened", methodName))
	}

	return result
}
