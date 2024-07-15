package double

import (
	"github.com/stretchr/testify/assert"
	"reflect"
)

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

type Arguments []interface{}

func (a Arguments) Matches(arguments ...interface{}) bool {
	if len(a) != len(arguments) {
		return false
	}

	for i, argument := range arguments {
		switch expected := a[i].(type) {
		case AnythingOfTypeArgument:
			if reflect.TypeOf(expected).Name() != string(expected) && reflect.TypeOf(argument).Name() != string(expected) {
				return false
			}
		case *IsTypeArgument:
			actualT := reflect.TypeOf(argument)
			if actualT != expected.t {
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
