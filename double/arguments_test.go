package double_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestArguments_Get(t *testing.T) {
	t.Run("Return the value of the given argument", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello", 123, true})

		assert.Equal(t, "hello", args.Get(0).(string))
		assert.Equal(t, 123, args.Get(1).(int))
		assert.Equal(t, true, args.Get(2).(bool))
	})

	t.Run("Panic if the argument does not exist", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello"})

		assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
			args.Get(1)
		})
	})
}

func TestArguments_String(t *testing.T) {
	t.Run("normal String() method - return a string representation of the args", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello", 123, true})

		assert.Equal(t, `string,int,bool`, args.String())
	})

	t.Run("Return the value of the given argument", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello", 123, true})

		assert.Equal(t, "hello", args.String(0))
	})

	t.Run("Panic if the argument does not exist", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello"})

		assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
			args.String(1)
		})
	})

	t.Run("Panic if the argument is not a string", func(t *testing.T) {
		var args = Arguments([]interface{}{true})

		assert.PanicsWithValue(t, "assert: arguments: String(0) failed because object wasn't correct type: %!s(bool=true)", func() {
			args.String(0)
		})
	})

	t.Run("Panic if we pass more than one index", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello", "world"})

		assert.PanicsWithValue(t, "assert: arguments: Wrong number of arguments passed to String.  Must be 0 or 1, not 2", func() {
			args.String(0, 1)
		})
	})
}

func TestArguments_Int(t *testing.T) {
	t.Run("Return the value of the given argument", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello", 123, true})

		assert.Equal(t, 123, args.Int(1))
	})

	t.Run("Panic if the argument does not exist", func(t *testing.T) {
		var args = Arguments([]interface{}{123})

		assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
			args.Int(1)
		})
	})

	t.Run("Panic if the argument is not an int", func(t *testing.T) {
		var args = Arguments([]interface{}{true})

		assert.PanicsWithValue(t, "assert: arguments: Int(0) failed because object wasn't correct type: %!s(bool=true)", func() {
			args.Int(0)
		})
	})
}

func TestArguments_Error(t *testing.T) {
	t.Run("Return the value of the given argument", func(t *testing.T) {
		var err = errors.New("an Error")
		var args = Arguments([]interface{}{"string", 123, true, err})

		assert.Equal(t, err, args.Error(3))
	})

	t.Run("Return nil of the given argument is nil", func(t *testing.T) {
		var args = Arguments([]interface{}{"string", 123, true, nil})

		assert.Equal(t, nil, args.Error(3))
	})

	t.Run("Panic if the argument does not exist", func(t *testing.T) {
		var args = Arguments([]interface{}{123})

		assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
			_ = args.Error(1)
		})
	})

	t.Run("Panic if the argument is not an error", func(t *testing.T) {
		var args = Arguments([]interface{}{true})

		assert.PanicsWithValue(t, "assert: arguments: Error(0) failed because object wasn't correct type: %!s(bool=true)", func() {
			_ = args.Error(0)
		})
	})
}

func TestArguments_Bool(t *testing.T) {
	t.Run("Return the value of the given argument", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello", 123, true})

		assert.Equal(t, true, args.Bool(2))
	})

	t.Run("Panic if the argument does not exist", func(t *testing.T) {
		var args = Arguments([]interface{}{"hello"})

		assert.PanicsWithValue(t, "assert: arguments: Cannot call Get(1) because there are 1 argument(s).", func() {
			args.Bool(1)
		})
	})

	t.Run("Panic if the argument is not a boolean", func(t *testing.T) {
		var args = Arguments([]interface{}{123})

		assert.PanicsWithValue(t, "assert: arguments: Bool(0) failed because object wasn't correct type: %!s(int=123)", func() {
			args.Bool(0)
		})
	})
}

func TestArguments_Matches(t *testing.T) {
	t.Run("false if length are different", func(t *testing.T) {
		var args = Arguments{123}
		assert.False(t, args.Matches(123, true))

		args = Arguments{123, true, 1.0}
		assert.False(t, args.Matches(123, true))
	})

	t.Run("compare primitives", func(t *testing.T) {
		var args = Arguments{123, true, 1.0}

		assert.True(t, args.Matches(123, true, 1.0))
		assert.False(t, args.Matches(456, true, 1.0))
		assert.False(t, args.Matches(123, false, 1.0))
		assert.False(t, args.Matches(123, true, 1.2))
	})

	t.Run("compare string", func(t *testing.T) {
		var args = Arguments{"String"}

		assert.True(t, args.Matches("String"))
		assert.False(t, args.Matches(""))
	})

	t.Run("compare list", func(t *testing.T) {
		var args = Arguments{[]int{1, 2, 3}}

		assert.True(t, args.Matches([]int{1, 2, 3}))
		assert.False(t, args.Matches([]int{4, 2, 3}))
	})

	t.Run("compare with Anything argument", func(t *testing.T) {
		var args = Arguments{1, Anything, 3}

		assert.True(t, args.Matches(1, 2, 3))
		assert.True(t, args.Matches(1, "String", 3))
		assert.True(t, args.Matches(1, 4.5, 3))
		assert.True(t, args.Matches(1, 2, Anything))
	})

	t.Run("compare with AnythingOfType argument", func(t *testing.T) {
		var args = Arguments{AnythingOfType("int"), AnythingOfType("string"), AnythingOfType("*double_test.ExampleType")}

		assert.True(t, args.Matches(1, "String", &ExampleType{}))
		assert.True(t, args.Matches(2, "any string", &ExampleType{}))
		assert.False(t, args.Matches("any string", "any string", &ExampleType{}))
		assert.False(t, args.Matches(2, 2, &ExampleType{}))
		assert.False(t, args.Matches(2, "any string", ExampleType{}))
	})

	t.Run("compare IsType argument", func(t *testing.T) {
		var args = Arguments([]interface{}{"string", IsType(0), true})

		assert.True(t, args.Matches("string", 123, true))
		assert.False(t, args.Matches("string", "string", true))
	})

	t.Run("compare Matcher argument", func(t *testing.T) {
		matchFn := func(a int) bool {
			return a == 123
		}
		var args = Arguments{"string", MatchedBy(matchFn), true}

		assert.True(t, args.Matches("string", 123, true))
		assert.False(t, args.Matches("string", 124, true))
		assert.False(t, args.Matches("string", false, true))
		assert.False(t, args.Matches("string", nil, true))
	})

	t.Run("compare Matcher argument with type nillable", func(t *testing.T) {
		matchFn := func(a []int) bool {
			return a == nil
		}
		var args = Arguments{"string", MatchedBy(matchFn), true}

		assert.True(t, args.Matches("string", nil, true))
		assert.False(t, args.Matches("string", 123, true))
		assert.False(t, args.Matches("string", false, true))
	})
}

func TestMatchedBy(t *testing.T) {
	t.Run("Create a ArgumentMatcher with a function with one argument and a boolean as return argument", func(t *testing.T) {
		matcher := MatchedBy(func(a int) bool { return a == 123 })

		assert.NotNil(t, matcher)
	})

	t.Run("Panic when argument is not a function", func(t *testing.T) {
		assert.PanicsWithValue(t, "assert: arguments: %!s(int=123) is not a func", func() {
			MatchedBy(123)
		})
	})

	t.Run("Panic when argument is a function with no argument", func(t *testing.T) {
		var fn interface{}
		fn = func() bool { return false }
		assert.PanicsWithValue(t, fmt.Sprintf("assert: arguments: %s does not take exactly one argument", fn), func() {
			MatchedBy(fn)
		})
	})

	t.Run("Panic when argument is a function with more than one argument", func(t *testing.T) {
		var fn interface{}
		fn = func(a int, b string) bool { return false }
		assert.PanicsWithValue(t, fmt.Sprintf("assert: arguments: %s does not take exactly one argument", fn), func() {
			MatchedBy(fn)
		})
	})

	t.Run("Panic when argument is a function with not a boolean as return argument", func(t *testing.T) {
		var fn interface{}
		fn = func(a int) {}
		assert.PanicsWithValue(t, fmt.Sprintf("assert: arguments: %s does not return a bool", fn), func() {
			MatchedBy(fn)
		})
	})

	t.Run("Panic when argument is a function with more than one boolean as return argument", func(t *testing.T) {
		var fn interface{}
		fn = func(a int) (bool, error) { return true, nil }
		assert.PanicsWithValue(t, fmt.Sprintf("assert: arguments: %s does not return a bool", fn), func() {
			MatchedBy(fn)
		})

	})
}
