package double

import "github.com/stretchr/testify/assert"

type Arguments []interface{}

func (a Arguments) Equal(arguments ...interface{}) bool {
	for i, argument := range arguments {
		if !assert.ObjectsAreEqual(a[i], argument) {
			return false
		}
	}
	return true
}
