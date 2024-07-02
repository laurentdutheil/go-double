package double_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/laurentdutheil/go-double/double"
)

func TestCall_Return_PredefineReturnArguments(t *testing.T) {
	call := NewCall("Method")

	call = call.Return(1, nil)

	assert.Equal(t, Arguments{1, nil}, call.ReturnArguments)
}
