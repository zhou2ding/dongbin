package rand

import (
	"blog/pkg/v"
	"crypto/md5"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"sync"
)

var gIdGenerator *IdGenerator

type IdGenerator struct {
	uuidVer int
	id      string
	mtx     sync.Mutex
}

func init() {
	uuidVer := v.GetViper().GetInt("uuid.v")
	gIdGenerator = &IdGenerator{
		uuidVer: uuidVer,
	}
}

func GetIdGeneratorInstance() *IdGenerator {
	return gIdGenerator
}

func (c *IdGenerator) GetId() string {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.id = c.GenId()
	return c.id
}

func (c *IdGenerator) GenId() string {
	var uid string
	switch c.uuidVer {
	case 1:
		uid = uuid.NewV1().String()
	case 2:
		uid = uuid.NewV2(uuid.DomainPerson).String()
	case 3:
		uid = uuid.NewV3(uuid.NamespaceDNS, "www.dongbin.com").String()

	case 4:
		uid = uuid.NewV4().String()
	case 5:
		uid = uuid.NewV5(uuid.NamespaceURL, "www.dongbin.com").String()
	}
	h5 := md5.New()
	h5.Write([]byte(uid))
	return hex.EncodeToString(h5.Sum(nil))
}
