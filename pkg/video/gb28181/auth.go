package gb28181

import (
	"crypto/md5"
	"fmt"
	"github.com/ghettovoice/gosip/sip"
)

func Verify(user, pass, realm, nonce string, auth *sip.Authorization) bool {
	s1 := fmt.Sprintf("%s:%s:%s", user, realm, pass)
	r1 := fmt.Sprintf("%x", md5.Sum([]byte(s1)))

	s2 := fmt.Sprintf("REGISTER:%s", auth.Uri())
	r2 := fmt.Sprintf("%x", md5.Sum([]byte(s2)))

	if r1 == "" || r2 == "" {
		return false
	}

	s3 := fmt.Sprintf("%s:%s:%s", r1, nonce, r2)
	r3 := fmt.Sprintf("%x", md5.Sum([]byte(s3)))

	return r3 == auth.Response()
}
