package double

import "testing"

// New is a constructor for Stub, Spy and Mock
//
//	type MockExample struct {
//		Mock
//	}
//	...
//	func TestExample(t *testing.T) {
//		myMock := New[MockExample](t)
//		...
//	}
func New[T any, TT tester[T]](t TestingT) *T {
	var result interface{} = new(T)
	tester := result.(TT)
	tester.Test(t)
	tester.Caller(result)
	return result.(*T)
}

type tester[T any] interface {
	Test(t TestingT)
	Caller(c interface{})
	*T
}

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
	Helper()
	FailNow()
}

// Check if TestingT interface can wrap testing.T
var _ TestingT = (*testing.T)(nil)
