package double

import (
	"fmt"
	"sync"
	"time"
)

// Call represents a method call and is used for providing answers on calls made during tests.
type Call struct {

	// The name of the method that will be called.
	MethodName string

	// Holds the arguments of the method.
	Arguments Arguments

	// Holds the arguments that should be returned when this method is called.
	ReturnArguments Arguments

	// The number of times to return the return arguments. 0 means to always return the values.
	times int

	// Amount of times this call has been called
	totalCalls int

	// panicMessage holds msg to be used to panic on the function call
	// if the panicMessage is set to a non nil string the function call will panic
	// irrespective of other settings
	panicMessage *string

	// Holds a channel that will be used to block the Return until it either
	// receives a message or is closed. nil means it returns immediately.
	waitFor  <-chan time.Time
	waitTime time.Duration

	// Holds a handler used to manipulate arguments content that are passed by
	// reference. It's useful when mocking methods such as unmarshalers or
	// decoders.
	runFn func(Arguments)

	mutex sync.Mutex
}

// NewCall constructor for Call
func NewCall(methodName string, arguments ...interface{}) *Call {
	return &Call{MethodName: methodName, Arguments: arguments}
}

// Return specifies the return arguments for the stubbed behaviour.
//
//	Stub.On("DoSomething").Return(errors.New("failed"))
func (c *Call) Return(arguments ...interface{}) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.ReturnArguments = append(c.ReturnArguments, arguments...)
	return c
}

// Once indicates that the mock should only return the value once.
//
//	Stub.On("Method", arg1, arg2).Return(returnArg1, returnArg2).Once()
func (c *Call) Once() *Call {
	return c.Times(1)
}

// Twice indicates that the mock should only return the value twice.
//
//	Stub.On("Method", arg1, arg2).Return(returnArg1, returnArg2).Twice()
func (c *Call) Twice() *Call {
	return c.Times(2)
}

// Times indicates that the mock should only return the indicated number of times.
//
//	Stub.On("Method", arg1, arg2).Return(returnArg1, returnArg2).Times(5)
func (c *Call) Times(i int) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.times = i
	return c
}

// Panic specifies if the function call should fail and the panic message
//
//	Stub.On("DoSomething").Panic("test panic")
func (c *Call) Panic(panicMessage string) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.panicMessage = &panicMessage
	return c
}

// WaitUntil sets the channel that will block the stub's return until its closed
// or a message is received.
//
//	Stub.On("Method", arg1, arg2).WaitUntil(time.After(time.Second))
func (c *Call) WaitUntil(w <-chan time.Time) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.waitFor = w
	return c
}

// After sets how long to block until the call returns
//
//	Stub.On("Method", arg1, arg2).After(time.Second)
func (c *Call) After(duration time.Duration) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.waitTime = duration
	return c
}

// Run sets a handler to be called before returning. It can be used when
// stubbing a method (such as an unmarshaler) that takes a pointer to a struct and
// sets properties in such struct
//
//	Stub.On("Unmarshal", AnythingOfType("*map[string]interface{}")).Return().Run(func(args Arguments) {
//		arg := args.Get(0).(*map[string]interface{})
//		arg["foo"] = "bar"
//	})
func (c *Call) Run(fn func(Arguments)) *Call {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.runFn = fn
	return c
}

// String return a string representation of a call
func (c *Call) String() string {
	return fmt.Sprintf("%s(%s)%s", c.MethodName, c.Arguments.String(), c.Arguments.valuesString())
}

func (c *Call) matches(t TestingT, methodName string, arguments ...interface{}) bool {
	return c.MethodName == methodName && c.Arguments.Matches(t, arguments...)
}

// canBeCalled return if the method call be called again
func (c *Call) canBeCalled() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.times == 0 || c.totalCalls < c.times
}

// calledPredefinedTimes return if the method was called the predefined times
func (c *Call) calledPredefinedTimes() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.times == 0 || c.times > 0 && c.times == c.totalCalls
}

// called executes the predefined behaviour of the call (waitFor, waitTime, panicMessage,,,)
// and return the predefined return arguments
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

// Calls collection of Call
type Calls []*Call

func (c *Calls) append(methodName string, arguments []interface{}) *Call {
	call := NewCall(methodName, arguments...)
	*c = append(*c, call)
	return call
}

// find the Call that matches methodName and arguments
// and check if the method can be called (Once, Twice, Times...)
// Return the null object noCallFound if no Call was found
func (c *Calls) find(t TestingT, methodName string, arguments ...interface{}) *Call {
	for _, predefinedCall := range *c {
		if predefinedCall.matches(t, methodName, arguments...) &&
			predefinedCall.canBeCalled() {
			return predefinedCall
		}
	}
	return noCallFound
}

var noCallFound = NewCall("-CallNotFound-")
