package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

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
