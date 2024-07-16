package double

type Call struct {
	MethodName      string
	Arguments       Arguments
	ReturnArguments Arguments
	times           int
	callCounter     int
	panicMessage    *string
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

func (c *Call) Twice() {
	c.times = 2
}

func (c *Call) Times(i int) {
	c.times = i
}

func (c *Call) alreadyCalledPredefinedTimes() bool {
	return c.times > 0 && c.times == c.callCounter
}

func (c *Call) incrementNumberOfCall() {
	c.callCounter++
}

func (c *Call) Panic(panicMessage string) {
	c.panicMessage = &panicMessage
}

type Method struct {
	Name   string
	NumOut int
}

type Calls []*Call

func (c Calls) find(methodName string, arguments ...interface{}) *Call {
	for _, predefinedCall := range c {
		if methodName == predefinedCall.MethodName {
			if !predefinedCall.Arguments.Matches(arguments...) || predefinedCall.alreadyCalledPredefinedTimes() {
				continue
			}
			return predefinedCall
		}
	}
	return noCallFound
}

var noCallFound = NewCall("-CallNotFound-")
