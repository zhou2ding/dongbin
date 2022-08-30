package rpcencrypt

import (
	"bytes"
	"encoding/base64"
	"sync"
)

const (
	Base64Length = 64
)

var (
	once                   sync.Once
	gEncryptHelperInstance *EncryptHelper
)

type EncryptHelper struct {
	encoding *base64.Encoding
	encoder  *encoder
}

func (helper *EncryptHelper) base64Encode(src []byte) string {
	return helper.encoding.EncodeToString(src)
}

func (helper *EncryptHelper) base64Decode(s string) []byte {
	quotient := len(s) / 4
	remainder := len(s) % 4
	var decLen int
	if remainder != 0 { //调整，使得待解密字符串长度为4的倍数 ，若余数为1直接丢弃，否则补充‘=’
		decLen = quotient*3 + remainder - 1
		switch remainder {
		case 1:
			s = s[:len(s)-1]
		case 2:
			s += "=="
		case 3:
			s += "="
		}
	} else {
		decLen = quotient * 3
	}

	if len(s) == 0 {
		return nil
	}

	decData, err := helper.encoding.DecodeString(s)
	if err != nil {
		return nil
	}

	return decData[:decLen]
}

func (helper *EncryptHelper) Encrypt(user string, pwd string, rand string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(pwd); i++ {
		if !helper.encoder.contains(pwd[i]) {
			buffer.WriteByte(helper.encoder.get(int(pwd[i]) % Base64Length))
		} else {
			buffer.WriteByte(pwd[i])
		}
	}
	pwd2 := helper.base64Decode(buffer.String())
	maxLen := len(pwd2)
	if len(user) > maxLen {
		maxLen = len(user)
	}
	if len(rand) > maxLen {
		maxLen = len(rand)
	}

	buffer.Reset()
	for idx := 0; idx < maxLen; idx++ {
		tmp := 0
		if idx < len(user) {
			tmp += int(user[idx])
		}
		if idx < len(pwd2) {
			tmp += int(pwd2[idx])
		}
		if idx < len(rand) {
			tmp += int(rand[idx])
		}
		tmp %= 256
		buffer.WriteByte(byte(tmp))
	}
	return helper.base64Encode(buffer.Bytes())
}

func GetEncryptHelperInstance() *EncryptHelper {
	once.Do(func() {
		ec := &encoder{
			codeStr: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/",
			verify: func(b byte) bool {
				if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '+' || b == '/' {
					return true
				}
				return false
			},
		}

		gEncryptHelperInstance = &EncryptHelper{
			encoding: base64.NewEncoding(ec.getEncoder()),
			encoder:  ec,
		}
	})

	return gEncryptHelperInstance
}
