package double

type Arguments []interface{}

func (a Arguments) Equal(arguments ...interface{}) bool {
	for i, argument := range arguments {
		if a[i] != argument {
			return false
		}
	}
	return true
}

type Call struct {
	MethodName      string
	Arguments       Arguments
	ReturnArguments Arguments
	times           int
	callCounter     int
}

func NewCall(methodName string, arguments ...interface{}) *Call {
	return &Call{MethodName: methodName, Arguments: arguments}
}

func (c *Call) Return(arguments ...interface{}) *Call {
	c.ReturnArguments = append(c.ReturnArguments, arguments...)
	return c
}

func (c *Call) Once() {
	c.times = 1
}

func (c *Call) alreadyCalledPredefinedTimes() bool {
	return c.times > 0 && c.times == c.callCounter
}

func (c *Call) updateNumberOfPredefinedCall() {
	if c.times > 0 {
		c.callCounter++
	}
}

type Method struct {
	Name   string
	NumOut int
}
