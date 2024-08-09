package double

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
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

func (a Arguments) Matches(arguments ...interface{}) bool {
	if len(a) != len(arguments) {
		return false
	}

	for i, actual := range arguments {
		expected := a[i]
		switch expectedType := expected.(type) {
		case AnythingOfTypeArgument:
			if reflect.TypeOf(actual).Name() != string(expectedType) && reflect.TypeOf(actual).String() != string(expectedType) {
				return false
			}
		case *IsTypeArgument:
			actualT := reflect.TypeOf(actual)
			if actualT != expectedType.t {
				return false
			}
		case ArgumentMatcher:
			matcher := expected.(ArgumentMatcher)
			var matches bool
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				matches = matcher.Matches(actual)
			}()
			if !matches {
				return false
			}
		default:
			if assert.ObjectsAreEqual(expected, Anything) || assert.ObjectsAreEqual(actual, Anything) {
				continue
			}
			if !assert.ObjectsAreEqual(expected, actual) {
				return false
			}
		}
	}
	return true
}

const Anything = "double.Anything"

type AnythingOfTypeArgument string

func AnythingOfType(t string) AnythingOfTypeArgument {
	return AnythingOfTypeArgument(t)
}

type IsTypeArgument struct {
	t reflect.Type
}

func IsType(t interface{}) *IsTypeArgument {
	return &IsTypeArgument{t: reflect.TypeOf(t)}
}

type ArgumentMatcher struct {
	fn reflect.Value
}

func (f ArgumentMatcher) Matches(argument interface{}) bool {
	expectType := f.fn.Type().In(0)

	argType := reflect.TypeOf(argument)
	var arg reflect.Value
	if argType == nil {
		arg = reflect.New(expectType).Elem()
	} else {
		arg = reflect.ValueOf(argument)
	}

	if argType == nil && !isNilSupported(expectType) {
		panic(errors.New("attempting to call matcher with nil for non-nil expected type"))
	}
	if argType == nil || argType.AssignableTo(expectType) {
		result := f.fn.Call([]reflect.Value{arg})
		return result[0].Bool()
	}
	return false
}

func isNilSupported(expectType reflect.Type) bool {
	switch expectType.Kind() {
	case reflect.Interface, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Ptr:
		return true
	default:
		return false
	}
}

func MatchedBy(fn interface{}) ArgumentMatcher {
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		panic(fmt.Sprintf("assert: arguments: %s is not a func", fn))
	}
	if fnType.NumIn() != 1 {
		panic(fmt.Sprintf("assert: arguments: %s does not take exactly one argument", fn))
	}
	if fnType.NumOut() != 1 || fnType.Out(0).Kind() != reflect.Bool {
		panic(fmt.Sprintf("assert: arguments: %s does not return a bool", fn))
	}

	return ArgumentMatcher{fn: reflect.ValueOf(fn)}
}
