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
func New[T any, DT Double[T]](t TestingT) *T {
	result := new(T)
	dt := DT(result)
	dt.Test(t)
	dt.Caller(result)
	return result
}

type Double[T any] interface {
	Test(t TestingT)
	Caller(caller interface{})
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
