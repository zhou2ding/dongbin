package rpcserver

import (
	"blog/pkg/internal/rpcpackage"
	"blog/pkg/l"
	"blog/pkg/rand"
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

const (
	RpcResponseBufferMax = 32
)

type rpcSession struct {
	user         *rpcUser
	sendChan     chan [][]byte
	rwMutex      sync.RWMutex
	id           int
	sessionId    int
	statusSendCh chan<- *SessionStatus
}

type SessionStatus struct {
	User      string
	SessionId int
	ClientId  int
	IsOnline  bool
	Ip        []byte
}

type jsonRequest struct {
	Domain string          `json:"domain"`
	Key    string          `json:"key"`
	Value  json.RawMessage `json:"value"`
	Id     uint32          `json:"id"`
}

type jsonResponse struct {
	Domain string          `json:"domain"`
	Id     uint32          `json:"id"`
	Key    string          `json:"key,omitempty"`
	Value  json.RawMessage `json:"value,omitempty"`
	RetVal int32           `json:"retval"`
	ErrMsg string          `json:"errmsg,omitempty"`
}

type jsonNotifyMessage struct {
	Domain string          `json:"domain"`
	Id     uint32          `json:"id"`
	Value  json.RawMessage `json:"value"`
}

func newRPCSession(checker UserChecker, sessionStatusReceiveCh chan<- *SessionStatus) *rpcSession {
	return &rpcSession{
		user:         newRpcUser(checker),
		sessionId:    rand.GetIdGeneratorInstance().GetId(),
		statusSendCh: sessionStatusReceiveCh,
	}
}

func (c *rpcSession) Id() int {
	return c.id
}

func (c *rpcSession) GetSessionId() int {
	return c.sessionId
}

func (c *rpcSession) OnTimeout() {
	l.GetLogger().Warn("RPCSession Timeout!!!", zap.Int("session Id", c.GetSessionId()))
	c.doClose()
}

func (c *rpcSession) Valid() bool {
	return c.user.getStatus() == SecondLogin
}

func (c *rpcSession) SendMessage(id int, value json.RawMessage) {
	msg := jsonNotifyMessage{
		Domain: "notify",
		Id:     uint32(id),
		Value:  value,
	}

	resp, _ := json.Marshal(msg)

	respPackets := rpcpackage.GetRPCMsgSplitterInstance().SplitPacket(resp, nil)

	c.doSend(respPackets)
}

func (c *rpcSession) Open() chan [][]byte {
	//向RPCSessionManager注册
	GetSessionMgr().register(c, time.Now().Unix())

	c.user.setStatus(NotLogin)
	c.sendChan = make(chan [][]byte, RpcResponseBufferMax)
	return c.sendChan
}

func (c *rpcSession) Close() {
	if c.doClose() == true {
		//从RPCSessionManager中注销
		GetSessionMgr().unregister(c)
	}
}

func (c *rpcSession) OnMessage(jsonReq []byte, binaryReq []byte) error {
	/*在这里可以对goroutine数量施加限制，从而限制并发的请求数,暂时不加限制*/
	//登录过程串行处理
	if c.user.getStatus() < SecondLogin {
		return c.doLogin(jsonReq, binaryReq)
	}

	//登录后的请求并发处理
	var jReq jsonRequest
	if err := json.Unmarshal(jsonReq, &jReq); err != nil {
		l.GetLogger().Warn("RPCSession OnMessage error: ", zap.Error(err), zap.Int("session Id", c.GetSessionId()))
		return err
	}

	//时间同步对处理的速度比较敏感，快速处理!!!
	//第一时间记录下请求到达服务器时的服务器unix时间戳
	if jReq.Key == "Systime.sync" {
		now := strconv.FormatInt(time.Now().UnixNano(), 16)
		binaryReq = []byte(now)
	}

	go c.doExecute(&jReq, binaryReq)

	return nil
}

func (c *rpcSession) doLogin(jsonReq []byte, binaryReq []byte) error {
	jsonResp, id, err := c.user.login(jsonReq)
	if err != nil {
		return err
	}

	if c.user.getStatus() == FirstLogin {
		c.id = id
		l.GetLogger().Info("RPCSession GetIdentifyInfo", zap.Int("id", c.id), zap.Int("session Id", c.GetSessionId()))
	}

	//二次登录完成后，更新下时间
	if c.user.getStatus() == SecondLogin {
		GetSessionMgr().updateAndDeleteInvalid(c, time.Now().Unix())

		l.GetLogger().Info("RPCSession doLogin set ListenState", zap.Int("Id", c.id), zap.Int("session Id", c.GetSessionId()))
		//此处由于项目需要ip这个字段，所以binaryReq里会装客户端的ip
		c.sendSessionState(true, binaryReq)
	}

	respPackets := rpcpackage.GetRPCMsgSplitterInstance().SplitPacket(jsonResp, nil)

	err = c.doSend(respPackets)
	if err != nil {
		l.GetLogger().Info("RPCSession doLogin doSend:", zap.Error(err), zap.Int("session Id", c.GetSessionId()))
	}

	return nil
}

func (c *rpcSession) doExecute(jReq *jsonRequest, binaryReq []byte) {
	if c.user.getStatus() == ClientLogout {
		//多请求并发处理，如果客户端登出，则未完成的请求不再继续处理
		return
	}

	var resp jsonResponse
	resp.Domain = jReq.Domain
	resp.Id = jReq.Id
	resp.Key = jReq.Key

	var binaryResp []byte
	var isLogout = false

	l.GetLogger().Info("RPCSession doExecute", zap.String("key", jReq.Key), zap.Int("session Id", c.GetSessionId()))

	//过滤订阅流程

	jsonResp, _ := json.Marshal(resp)

	respPackets := rpcpackage.GetRPCMsgSplitterInstance().SplitPacket(jsonResp, binaryResp)

	bt := time.Now()
	err := c.doSend(respPackets)
	if err != nil {
		l.GetLogger().Info("RPCSession doExecute doSend", zap.Error(err), zap.Int("session Id", c.GetSessionId()))
	}
	et := time.Since(bt)
	l.GetLogger().Info("RPCSession doExecute doSend used time", zap.Int64("time", int64(et)/int64(time.Millisecond)), zap.Int("session Id", c.GetSessionId()))

	if isLogout { //响应客户端登出请求后关闭会话
		c.Close()
	}
}

func (c *rpcSession) doSend(packets [][]byte) error {
	c.rwMutex.RLock()
	defer c.rwMutex.RUnlock()
	if c.sendChan != nil {
		c.sendChan <- packets
	} else {
		return errors.New("sendChan closed")
	}

	return nil
}

func (c *rpcSession) sendSessionState(state bool, ip []byte) {
	l.GetLogger().Info("RPCSession sendSessionState", zap.Bool("state", state), zap.Int("session Id", c.GetSessionId()))

	if !state { //没有登录成功就退出，不设置状态
		if c.user.getStatus() < SecondLogin {
			return
		}
	}

	statusMsg := SessionStatus{
		User:      c.user.user,
		SessionId: c.GetSessionId(),
		ClientId:  c.id,
		IsOnline:  state,
		Ip:        ip,
	}
	c.statusSendCh <- &statusMsg
}

func (c *rpcSession) doClose() bool {
	isCloseDone := false

	c.rwMutex.Lock()
	if c.sendChan != nil {
		close(c.sendChan)
		c.sendChan = nil

		isCloseDone = true
	}
	c.rwMutex.Unlock()

	if isCloseDone {
		if c.user.getStatus() == SecondLogin {
			c.user.setStatus(ClientLogout)
		}

		l.GetLogger().Info("RPCSession doClose set ListenState", zap.Int("Id", c.id), zap.Int("session Id", c.GetSessionId()))
		c.sendSessionState(false, nil)

		//客户端有可能未注销订阅就（因异常等）断开了
	}

	return isCloseDone
}
