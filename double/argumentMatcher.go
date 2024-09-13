package double

import (
	"fmt"
	"reflect"
)

type ArgumentMatcher interface {
	matches(actual interface{}) bool
}

const Anything = "double.Anything"

func AnythingOfType(t string) ArgumentMatcher {
	return anythingOfTypeArgument(t)
}

func IsType(t interface{}) ArgumentMatcher {
	return &isTypeArgument{t: reflect.TypeOf(t)}
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

	return functionMatcherArgument{fn: reflect.ValueOf(fn)}
}

type anythingOfTypeArgument string

func (t anythingOfTypeArgument) matches(actual interface{}) bool {
	return reflect.TypeOf(actual).Name() == string(t) || reflect.TypeOf(actual).String() == string(t)
}

type isTypeArgument struct {
	t reflect.Type
}

func (t isTypeArgument) String() string {
	return t.t.Name()
}

func (t isTypeArgument) matches(actual interface{}) bool {
	return reflect.TypeOf(actual) == t.t
}

type functionMatcherArgument struct {
	fn reflect.Value
}

func (f functionMatcherArgument) String() string {
	return fmt.Sprintf("func(%s) bool", f.fn.Type().In(0).String())
}

func (f functionMatcherArgument) matches(argument interface{}) bool {
	expectType := f.fn.Type().In(0)

	argType := reflect.TypeOf(argument)
	var arg reflect.Value
	if argType == nil {
		arg = reflect.New(expectType).Elem()
	} else {
		arg = reflect.ValueOf(argument)
	}

	if argType == nil && !isNilSupported(expectType) {
		return false
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
