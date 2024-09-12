package double_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestArgumentMatcher(t *testing.T) {
	t.Run("MatchedBy", func(t *testing.T) {
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
	})
}
