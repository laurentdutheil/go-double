package double

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
)

type Arguments []interface{}

func (a Arguments) Matches(arguments ...interface{}) bool {
	if len(a) != len(arguments) {
		return false
	}

	for i, argument := range arguments {
		switch expectedType := a[i].(type) {
		case AnythingOfTypeArgument:
			if reflect.TypeOf(expectedType).Name() != string(expectedType) && reflect.TypeOf(argument).Name() != string(expectedType) {
				return false
			}
		case *IsTypeArgument:
			actualT := reflect.TypeOf(argument)
			if actualT != expectedType.t {
				return false
			}
		case ArgumentMatcher:
			matcher := a[i].(ArgumentMatcher)
			var matches bool
			func() {
				defer func() {
					if r := recover(); r != nil {
					}
				}()
				matches = matcher.Matches(argument)
			}()
			if !matches {
				return false
			}
		default:
			if assert.ObjectsAreEqual(a[i], Anything) || assert.ObjectsAreEqual(argument, Anything) {
				continue
			}
			if !assert.ObjectsAreEqual(a[i], argument) {
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
