package double

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
)

type Arguments []interface{}

// Get Returns the argument at the specified index.
func (a Arguments) Get(index int) interface{} {
	if index+1 > len(a) {
		panic(fmt.Sprintf("assert: arguments: Cannot call Get(%d) because there are %d argument(s).", index, len(a)))
	}
	return a[index]
}

// String gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
//
// If no index is provided, String() returns a complete string representation
// of the arguments.
func (a Arguments) String(indexOrNil ...int) string {
	if len(indexOrNil) == 0 {
		// normal String() method - return a string representation of the args
		var argsStr []string
		for _, arg := range a {
			argsStr = append(argsStr, fmt.Sprintf("%T", arg)) // handles nil nicely
		}
		return strings.Join(argsStr, ",")
	} else if len(indexOrNil) == 1 {
		// Index has been specified - get the argument at that index
		index := indexOrNil[0]
		var s string
		var ok bool
		if s, ok = a.Get(index).(string); !ok {
			panic(fmt.Sprintf("assert: arguments: String(%d) failed because object wasn't correct type: %s", index, a.Get(index)))
		}
		return s
	}

	panic(fmt.Sprintf("assert: arguments: Wrong number of arguments passed to String.  Must be 0 or 1, not %d", len(indexOrNil)))
}

func (a Arguments) valuesString() string {
	if len(a) == 0 {
		return ""
	}

	var argVals []string
	for argIndex, arg := range a {
		argVals = append(argVals, fmt.Sprintf("%d: %#v", argIndex, arg))
	}
	return fmt.Sprintf("\n\t\t%s", strings.Join(argVals, "\n\t\t"))

}

// Error gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (a Arguments) Error(index int) error {
	obj := a.Get(index)
	var s error
	var ok bool
	if obj == nil {
		return nil
	}
	if s, ok = obj.(error); !ok {
		panic(fmt.Sprintf("assert: arguments: Error(%d) failed because object wasn't correct type: %s", index, obj))
	}
	return s
}

// Int gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (a Arguments) Int(index int) int {
	var s int
	var ok bool
	if s, ok = a.Get(index).(int); !ok {
		panic(fmt.Sprintf("assert: arguments: Int(%d) failed because object wasn't correct type: %s", index, a.Get(index)))
	}
	return s
}

// Bool gets the argument at the specified index. Panics if there is no argument, or
// if the argument is of the wrong type.
func (a Arguments) Bool(index int) bool {
	var s bool
	var ok bool
	if s, ok = a.Get(index).(bool); !ok {
		panic(fmt.Sprintf("assert: arguments: Bool(%d) failed because object wasn't correct type: %s", index, a.Get(index)))
	}
	return s
}
func (a Arguments) Matches(t TestingT, arguments ...interface{}) bool {
	t.Helper()

	if len(a) != len(arguments) {
		t.Logf("Arguments have not the same size: len(expected) == %d ; len(actual) == %d", len(a), len(arguments))
		return false
	}

	result := true
	for i, actual := range arguments {
		actualFmt := fmt.Sprintf("(%[1]T=%[1]v)", actual)
		expected := a[i]
		expectedFmt := fmt.Sprintf("(%[1]T=%[1]v)", expected)
		expectedType, ok := expected.(ArgumentMatcher)
		if ok {
			objectsMatch := expectedType.matches(actual)
			if objectsMatch {
				t.Logf("\t%d: PASS: %s matches %s", i, actualFmt, expectedFmt)
			} else {
				t.Logf("\t%d: FAIL: %s doesn't match %s", i, actualFmt, expectedFmt)
			}
			result = result && objectsMatch
		} else {
			objectAreEqual := assert.ObjectsAreEqual(expected, Anything) || assert.ObjectsAreEqual(actual, Anything) || assert.ObjectsAreEqual(expected, actual)
			if objectAreEqual {
				t.Logf("\t%d: PASS: %s == %s", i, actualFmt, expectedFmt)
			} else {
				t.Logf("\t%d: FAIL: %s != %s", i, actualFmt, expectedFmt)
			}
			result = result && objectAreEqual
		}
	}

	return result
}
