package double_test

import (
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
		var args = Arguments{AnythingOfType("int"), AnythingOfType("string"), 1.0}

		assert.True(t, args.Matches(1, "String", 1.0))
		assert.True(t, args.Matches(2, "any string", 1.0))
		assert.False(t, args.Matches("any string", "any string", 1.0))
		assert.False(t, args.Matches(2, 2, 1.0))
	})

	t.Run("compare IsType argument", func(t *testing.T) {
		var args = Arguments([]interface{}{"string", IsType(0), true})

		assert.True(t, args.Matches("string", 123, true))
		assert.False(t, args.Matches("string", "string", true))
	})

}
