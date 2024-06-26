package double

type Call struct {
	MethodName string
	Arguments  Arguments
}

type Arguments []interface{}
