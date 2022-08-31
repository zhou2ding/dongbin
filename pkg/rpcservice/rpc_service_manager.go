package rpcservice

import (
	"blog/pkg/l"
	"github.com/pkg/errors"
	"strings"
)

type RpcService interface {
	RPCCall(key string, jsonReq []byte, binaryReq []byte, extraInfo interface{}) *RPCReply
}

type RPCReply struct {
	RetVal         int32
	ErrMsg         string
	Value          []byte
	BinaryResponse []byte
}

func GetRPCServiceMgr() *RPCServiceManager {
	return gRPCServiceManager
}

var gRPCServiceManager = &RPCServiceManager{
	serviceMap: make(map[string]RpcService),
}

type RPCServiceManager struct {
	serviceMap map[string]RpcService
	notifyCh   chan interface{}
}

func (c *RPCServiceManager) Register(serviceName string, service RpcService) {
	l.GetLogger().Info("register " + serviceName)
	c.serviceMap[serviceName] = service
}

func (c *RPCServiceManager) Unregister(serviceName string) {
	delete(c.serviceMap, serviceName)
}

func (c *RPCServiceManager) Execute(key string, jsonReq []byte, binaryReq []byte, extraInfo interface{}) *RPCReply {
	idx := strings.IndexByte(key, '.')
	if idx < 0 {
		l.GetLogger().Warn("RPCMethodManager Execute . Not Found")
		return &RPCReply{
			RetVal: ErrUnknownReqType,
			ErrMsg: StatusText(ErrUnknownReqType),
		}
	}

	serviceName := key[0:idx]

	srvc, ok := c.serviceMap[serviceName]
	if !ok {
		l.GetLogger().Warn("RPCMethodManager Execute service Not Found")
		return &RPCReply{
			RetVal: ErrServiceNotSupport,
			ErrMsg: StatusText(ErrServiceNotSupport),
		}
	}

	return srvc.RPCCall(key, jsonReq, binaryReq, extraInfo)
}

func (c *RPCServiceManager) FilterAndPreprocess(key string) (string, string) {
	keyFields := strings.Split(key, ".")
	if len(keyFields) != 2 { //service.method
		return "none", "none"
	}

	//先判断service, 再判断method
	_, ok := c.serviceMap[keyFields[0]]
	if !ok {
		return "none", "none"
	}

	if keyFields[1] == "attach" || keyFields[1] == "detach" {
		return keyFields[0], keyFields[1]
	}

	return "none", "none"
}

func (c *RPCServiceManager) SetNotifyChan(ch chan interface{}) {
	c.notifyCh = ch
}

func (c *RPCServiceManager) Notify(msg *RPCMsg) {
	c.notifyCh <- msg
}

func (c *RPCServiceManager) StartTopicListen(key string) error {
	l.GetLogger().Info("StartTopicListen")
	srvc := c.serviceMap[key]
	reply := srvc.RPCCall(key+".attach", nil, nil, nil)
	if reply.RetVal != 0 {
		return errors.New(reply.ErrMsg)
	}

	return nil
}

func (c *RPCServiceManager) StopTopicListen(key string) error {
	l.GetLogger().Info("StopTopicListen")
	srvc := c.serviceMap[key]
	reply := srvc.RPCCall(key+".detach", nil, nil, nil)
	if reply.RetVal != 0 {
		return errors.New(reply.ErrMsg)
	}

	return nil
}
