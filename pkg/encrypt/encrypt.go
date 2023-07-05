package encrypt

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func Encrypt(s string, salts ...string) string {
	m5 := md5.New()
	m5.Write([]byte(s))
	if l := len(salts); l > 0 {
		arr := make([]string, l+1)
		m5.Write([]byte(fmt.Sprintf(strings.Join(arr, "%v"), salts)))
	}
	return hex.EncodeToString(m5.Sum(nil))
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
