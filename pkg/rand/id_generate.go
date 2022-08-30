package rand

import "sync"

var gIdGenerator *IdGenerator = &IdGenerator{id: 0}

type IdGenerator struct {
	id  int
	mtx sync.Mutex
}

func GetIdGeneratorInstance() *IdGenerator {
	return gIdGenerator
}

func (c *IdGenerator) GetId() int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.id++
	return c.id
}
