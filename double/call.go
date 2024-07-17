package double

import (
	"time"
)

type Call struct {
	MethodName      string
	Arguments       Arguments
	ReturnArguments Arguments
	times           int
	totalCalls      int
	panicMessage    *string
	waitFor         <-chan time.Time
	waitTime        time.Duration
	runFn           func(Arguments)
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

func (c *Call) Panic(panicMessage string) {
	c.panicMessage = &panicMessage
}

func (c *Call) WaitUntil(w <-chan time.Time) {
	c.waitFor = w
}

func (c *Call) After(duration time.Duration) {
	c.waitTime = duration
}

func (c *Call) Run(fn func(Arguments)) {
	c.runFn = fn
}

func (c *Call) alreadyCalledPredefinedTimes() bool {
	return c.times > 0 && c.times == c.totalCalls
}

func (c *Call) called(arguments ...interface{}) Arguments {
	c.totalCalls++

	if c.waitFor != nil {
		<-c.waitFor
	} else {
		time.Sleep(c.waitTime)
	}

	if c.panicMessage != nil {
		panic(*c.panicMessage)
	}

	if c.runFn != nil {
		c.runFn(arguments)
	}

	return c.ReturnArguments
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
