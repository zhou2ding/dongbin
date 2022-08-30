package rpcserver

import (
	"blog/pkg/l"
	"go.uber.org/zap"
	"time"
)

const (
	SessionCheckTimeCell = 5
)

type sessionGuarder interface {
	Id() int
	GetSessionId() int
	OnTimeout()
}

type sessionOperateType int

const (
	sessionRegister sessionOperateType = iota
	sessionUpdate
	sessionUpdateAndClear
	sessionUnregister
)

type sessionOperateMessage struct {
	guarder     sessionGuarder
	operateType sessionOperateType
	operateTime int64
}

type sessionGuardCtx struct {
	guarder     sessionGuarder
	lastOptTime int64
}

var gRPCSessionManager = &rpcSessionManager{
	guardContexts: make(map[int]*sessionGuardCtx),
	messageCh:     make(chan *sessionOperateMessage, 10),
}

func GetRPCSessionMgrInstance() *rpcSessionManager {
	return gRPCSessionManager
}

type rpcSessionManager struct {
	guardContexts map[int]*sessionGuardCtx
	messageCh     chan *sessionOperateMessage
}

func (c *rpcSessionManager) Start(aliveTime int64) {
	go c.keepSessionGuarded(aliveTime)
}

func (c *rpcSessionManager) register(guarder sessionGuarder, operateTime int64) {
	c.messageCh <- &sessionOperateMessage{
		guarder:     guarder,
		operateType: sessionRegister,
		operateTime: operateTime,
	}
}

func (c *rpcSessionManager) unregister(guarder sessionGuarder) {
	c.messageCh <- &sessionOperateMessage{
		guarder:     guarder,
		operateType: sessionUnregister,
	}
}

func (c *rpcSessionManager) update(guarder sessionGuarder, operateTime int64) {
	c.messageCh <- &sessionOperateMessage{
		guarder:     guarder,
		operateType: sessionUpdate,
		operateTime: operateTime,
	}
}

func (c *rpcSessionManager) updateAndDeleteInvalid(guarder sessionGuarder, operateTime int64) {
	c.messageCh <- &sessionOperateMessage{
		guarder:     guarder,
		operateType: sessionUpdateAndClear,
		operateTime: operateTime,
	}
}

func (c *rpcSessionManager) keepSessionGuarded(aliveTime int64) {
	ticker := time.NewTicker(time.Duration(SessionCheckTimeCell) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			now := t.Unix()
			l.GetLogger().Info("keepSessionGuarded", zap.Time("now", t))
			for id, ctx := range c.guardContexts {
				if now > ctx.lastOptTime && now-ctx.lastOptTime > aliveTime {
					ctx.guarder.OnTimeout()
					delete(c.guardContexts, id)
				}
			}
		case msg := <-c.messageCh:
			c.parseSessionOperateMessage(msg)
		}
	}
}

func (c *rpcSessionManager) parseSessionOperateMessage(msg *sessionOperateMessage) {
	switch msg.operateType {
	case sessionRegister:
		l.GetLogger().Info("rpcSessionManager register ", zap.Int("session Id", msg.guarder.GetSessionId()))
		_, ok := c.guardContexts[msg.guarder.GetSessionId()]
		if ok {
			l.GetLogger().Warn("repeat register sessionMonitor!")
			return
		}
		c.guardContexts[msg.guarder.GetSessionId()] = &sessionGuardCtx{
			guarder:     msg.guarder,
			lastOptTime: msg.operateTime,
		}
	case sessionUpdate:
		fallthrough
	case sessionUpdateAndClear:
		l.GetLogger().Info("rpcSessionManager update", zap.Int("session Id", msg.guarder.GetSessionId()))
		_, ok := c.guardContexts[msg.guarder.GetSessionId()]
		if ok {
			c.guardContexts[msg.guarder.GetSessionId()].lastOptTime = msg.operateTime
		} else {
			l.GetLogger().Warn("update non-existent sessionGuarder!")
			return
		}

		if msg.operateType == sessionUpdateAndClear {
			//相同的用户，新会话建立时，若有老会话存在，删除老会话
			for sessionId, ctx := range c.guardContexts {
				if ctx.guarder.Id() == msg.guarder.Id() && sessionId != msg.guarder.GetSessionId() {
					l.GetLogger().Info("rpcSessionManager delete old sessionGuarder", zap.Int("session Id", sessionId))
					ctx.guarder.OnTimeout()
					delete(c.guardContexts, sessionId)
				}
			}
		}
	case sessionUnregister:
		l.GetLogger().Info("rpcSessionManager unregister ", zap.Int("session Id", msg.guarder.GetSessionId()))
		_, ok := c.guardContexts[msg.guarder.GetSessionId()]
		if ok {
			delete(c.guardContexts, msg.guarder.GetSessionId())
		} else {
			l.GetLogger().Warn("unregister non-existent sessionGuarder!")
		}
	}
}
