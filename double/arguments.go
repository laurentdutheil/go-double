package double

import "github.com/stretchr/testify/assert"

type Arguments []interface{}

func (a Arguments) Equal(arguments ...interface{}) bool {
	if len(a) != len(arguments) {
		return false
	}

	for i, argument := range arguments {
		if !assert.ObjectsAreEqual(a[i], argument) {
			return false
		}
	}
	return true
}
