package rand

import (
	"math/rand"
	"sync"
	"time"
)

var once sync.Once
var gRandomGenerator *RandomGenerator

type RandomGenerator struct {
	rand      *rand.Rand
	sourceStr string
}

func (c *RandomGenerator) init() {
	c.rand = rand.New(rand.NewSource(time.Now().Unix()))
	c.sourceStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
}

func (c *RandomGenerator) GetRandomString(l int) string {
	r := make([]byte, 0, l)
	for i := 0; i < l; i++ {
		r = append(r, c.sourceStr[c.rand.Intn(len(c.sourceStr))])
	}
	return string(r)
}

func GetRandGeneratorInstance() *RandomGenerator {
	once.Do(func() {
		gRandomGenerator = &RandomGenerator{}
		gRandomGenerator.init()
	})
	return gRandomGenerator
}
