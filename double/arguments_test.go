package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestArguments_Equals(t *testing.T) {
	t.Run("false if length are different", func(t *testing.T) {
		var args = Arguments{123}
		assert.False(t, args.Equal(123, true))

		args = Arguments{123, true, 1.0}
		assert.False(t, args.Equal(123, true))
	})

	t.Run("compare primitives", func(t *testing.T) {
		var args = Arguments{123, true, 1.0}

		assert.True(t, args.Equal(123, true, 1.0))
		assert.False(t, args.Equal(456, true, 1.0))
		assert.False(t, args.Equal(123, false, 1.0))
		assert.False(t, args.Equal(123, true, 1.2))
	})

	t.Run("compare string", func(t *testing.T) {
		var args = Arguments{"String"}

		assert.True(t, args.Equal("String"))
		assert.False(t, args.Equal(""))
	})

	t.Run("compare list", func(t *testing.T) {
		var args = Arguments{[]int{1, 2, 3}}

		assert.True(t, args.Equal([]int{1, 2, 3}))
		assert.False(t, args.Equal([]int{4, 2, 3}))
	})
}
