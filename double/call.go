package double

import (
	"fmt"
	"sync"
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
	mutex           sync.Mutex
}

func NewCall(methodName string, arguments ...interface{}) *Call {
	return &Call{MethodName: methodName, Arguments: arguments}
}

func (c *Call) Return(arguments ...interface{}) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.ReturnArguments = append(c.ReturnArguments, arguments...)
	return c
}

func (c *Call) Once() *Call {
	return c.Times(1)
}

func (c *Call) Twice() *Call {
	return c.Times(2)
}

func (c *Call) Times(i int) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.times = i
	return c
}

func (c *Call) Panic(panicMessage string) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.panicMessage = &panicMessage
	return c
}

func (c *Call) WaitUntil(w <-chan time.Time) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.waitFor = w
	return c
}

func (c *Call) After(duration time.Duration) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.waitTime = duration
	return c
}

func (c *Call) Run(fn func(Arguments)) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.runFn = fn
	return c
}

func (c *Call) String() string {
	return fmt.Sprintf("%s(%s)%s", c.MethodName, c.Arguments.String(), c.Arguments.valuesString())
}

func (c *Call) canBeCalled() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.times == 0 || c.totalCalls < c.times
}

func (c *Call) calledPredefinedTimes() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.times == 0 || c.times > 0 && c.times == c.totalCalls
}

func (c *Call) called(arguments ...interface{}) Arguments {
	c.mutex.Lock()
	defer c.mutex.Unlock()
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
			if predefinedCall.Arguments.Matches(arguments...) && predefinedCall.canBeCalled() {
				return predefinedCall
			}
		}
	}
	return noCallFound
}

var noCallFound = NewCall("-CallNotFound-")
