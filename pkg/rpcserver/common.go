package rpcserver

import "sync"

type atomicStatus struct {
	value int
	mtx   sync.RWMutex
}

func (c *atomicStatus) GetValue() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return c.value
}

func (c *atomicStatus) SetValue(val int) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.value = val
}
