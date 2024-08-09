package double

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCall_String(t *testing.T) {
	t.Run("print method name", func(t *testing.T) {
		call := NewCall("Method")

		assert.Equal(t, "Method()", call.String())
	})

	t.Run("print method name and arguments", func(t *testing.T) {
		call := NewCall("Method", 1, "2", 3.4)

		assert.Contains(t, call.String(), "Method(int,string,float64)")
	})

	t.Run("print method name, arguments and values", func(t *testing.T) {
		call := NewCall("Method", 1, "2", nil)

		expected := `Method(int,string,<nil>)
		0: 1
		1: "2"
		2: <nil>`
		assert.Equal(t, expected, call.String())
	})
}
