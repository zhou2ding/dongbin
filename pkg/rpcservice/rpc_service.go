package rpcservice

type Method func(jsonReq []byte, binaryReq []byte, extraInfo interface{}) *RPCReply

type Base struct {
	MethodMap map[string]Method
}

func (c *Base) AddMethod(name string, method Method) {
	c.MethodMap[name] = method
}

func (c *Base) RemoveMethod(name string) {
	delete(c.MethodMap, name)
}

func (c *Base) RPCCall(key string, jsonReq []byte, binaryReq []byte, extraInfo interface{}) *RPCReply {
	if fn, ok := c.MethodMap[key]; ok {
		return fn(jsonReq, binaryReq, extraInfo)
	}

	reply := &RPCReply{
		RetVal: ErrMethodNotImplement,
		ErrMsg: StatusText(ErrMethodNotImplement),
	}

	return reply
}
