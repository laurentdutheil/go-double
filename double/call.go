package double

type Arguments []interface{}

type Call struct {
	MethodName      string
	Arguments       Arguments
	ReturnArguments Arguments
}

func NewCall(methodName string, arguments ...interface{}) *Call {
	return &Call{MethodName: methodName, Arguments: arguments}
}

func (c *Call) Return(arguments ...interface{}) *Call {
	c.ReturnArguments = append(c.ReturnArguments, arguments...)
	return c
}

type Method struct {
	Name   string
	NumOut int
}
