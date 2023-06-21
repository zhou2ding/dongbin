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

type action int

const (
	register action = iota
	update
	updateAndClear
	unregister
)

type sessionOperateMessage struct {
	guarder sessionGuarder
	action  action
	time    int64
}

type sessionGuardCtx struct {
	guarder     sessionGuarder
	lastOptTime int64
}

var gSessionMgr = &sessionMgr{
	guardContexts: make(map[int]*sessionGuardCtx),
	messageCh:     make(chan *sessionOperateMessage, 10),
}

func GetSessionMgr() *sessionMgr {
	return gSessionMgr
}

type sessionMgr struct {
	guardContexts map[int]*sessionGuardCtx
	messageCh     chan *sessionOperateMessage
}

func (c *sessionMgr) Start(aliveTime int64) {
	go c.keepSessionGuarded(aliveTime)
}

func (c *sessionMgr) register(guarder sessionGuarder, operateTime int64) {
	c.messageCh <- &sessionOperateMessage{
		guarder: guarder,
		action:  register,
		time:    operateTime,
	}
}

func (c *sessionMgr) unregister(guarder sessionGuarder) {
	c.messageCh <- &sessionOperateMessage{
		guarder: guarder,
		action:  unregister,
	}
}

func (c *sessionMgr) update(guarder sessionGuarder, operateTime int64) {
	c.messageCh <- &sessionOperateMessage{
		guarder: guarder,
		action:  update,
		time:    operateTime,
	}
}

func (c *sessionMgr) updateAndDeleteInvalid(guarder sessionGuarder, operateTime int64) {
	c.messageCh <- &sessionOperateMessage{
		guarder: guarder,
		action:  updateAndClear,
		time:    operateTime,
	}
}

func (c *sessionMgr) keepSessionGuarded(aliveTime int64) {
	ticker := time.NewTicker(time.Duration(SessionCheckTimeCell) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			now := t.Unix()
			l.Logger().Info("keepSessionGuarded", zap.Time("now", t))
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

func (c *sessionMgr) parseSessionOperateMessage(msg *sessionOperateMessage) {
	switch msg.action {
	case register:
		l.Logger().Info("sessionMgr register ", zap.Int("session Id", msg.guarder.GetSessionId()))
		_, ok := c.guardContexts[msg.guarder.GetSessionId()]
		if ok {
			l.Logger().Warn("repeat register sessionMonitor!")
			return
		}
		c.guardContexts[msg.guarder.GetSessionId()] = &sessionGuardCtx{
			guarder:     msg.guarder,
			lastOptTime: msg.time,
		}
	case update:
		fallthrough
	case updateAndClear:
		l.Logger().Info("sessionMgr update", zap.Int("session Id", msg.guarder.GetSessionId()))
		_, ok := c.guardContexts[msg.guarder.GetSessionId()]
		if ok {
			c.guardContexts[msg.guarder.GetSessionId()].lastOptTime = msg.time
		} else {
			l.Logger().Warn("update non-existent sessionGuarder!")
			return
		}

		if msg.action == updateAndClear {
			//相同的用户，新会话建立时，若有老会话存在，删除老会话
			for sessionId, ctx := range c.guardContexts {
				if ctx.guarder.Id() == msg.guarder.Id() && sessionId != msg.guarder.GetSessionId() {
					l.Logger().Info("sessionMgr delete old sessionGuarder", zap.Int("session Id", sessionId))
					ctx.guarder.OnTimeout()
					delete(c.guardContexts, sessionId)
				}
			}
		}
	case unregister:
		l.Logger().Info("sessionMgr unregister ", zap.Int("session Id", msg.guarder.GetSessionId()))
		_, ok := c.guardContexts[msg.guarder.GetSessionId()]
		if ok {
			delete(c.guardContexts, msg.guarder.GetSessionId())
		} else {
			l.Logger().Warn("unregister non-existent sessionGuarder!")
		}
	}
}
