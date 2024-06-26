package double

type Arguments []interface{}

type Call struct {
	MethodName      string
	Arguments       Arguments
	ReturnArguments Arguments
}

func NewCall(methodName string, arguments Arguments) *Call {
	return &Call{MethodName: methodName, Arguments: arguments}
}

func (c *Call) Return(args ...interface{}) *Call {
	c.ReturnArguments = append(c.ReturnArguments, args...)
	return c
}
